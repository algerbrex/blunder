package engine

const (
	SearchTTSize       = 16
	PerftTTSize        = 10
	TTEntrySizeInBytes = 16
	BucketSize         = 4

	LeftMost56Bits uint64 = 0xffffffffffffff00
	RightMost8Bits uint64 = 0xff

	LeftMost2Bits  uint8 = 0xc0
	RightMost6Bits uint8 = 0x3f
)

const (
	FailLowNode uint8 = iota
	FailHighNode
	PVNode
)

type SearchEntry struct {
	Hash           uint64
	BestMove       uint32
	Depth          int8
	Score          int16
	NodeTypeAndAge uint8
}

func (entry *SearchEntry) GetScoreAndBestMove(hash uint64, ply uint8, depth int8, alpha, beta int16) (int16, uint32, bool) {
	shouldUse := false
	score := int16(0)
	bestMove := NullMove

	if entry.Hash == hash {
		bestMove = entry.BestMove

		if entry.Depth >= depth {
			entryScore := entry.Score
			nodeType := (entry.NodeTypeAndAge & LeftMost2Bits) >> 6

			// If the score we get from the transposition table is a checkmate score, we need
			// to do a little extra work. This is because we store checkmates in the table using
			// their distance from the node they're found in, not their distance from the root.
			// So if we found a checkmate-in-8 in a node that was 5 plies from the root, we need
			// to store the score as a checkmate-in-3. Then, if we read the checkmate-in-3 from
			// the table in a node that's 4 plies from the root, we need to return the score as
			// checkmate-in-7.

			if entryScore > CheckmateThreshold {
				entryScore -= int16(ply)
			} else if entryScore < -CheckmateThreshold {
				entryScore += int16(ply)
			}

			if nodeType == PVNode {
				// If we have an exact entry, we can use the saved score.
				score = entryScore
				shouldUse = true
			} else if nodeType == FailLowNode && entryScore <= alpha {
				// If we have an alpha entry, and the entry's score is less than our
				// current alpha, then we know that our current alpha is probably the
				// best score we can get in this node, so we can stop searching and use
				// alpha.
				score = alpha
				shouldUse = true
			} else if nodeType == FailHighNode && entryScore >= beta {
				// If we have a beta entry, and the entry's score is greater than our
				// current beta, then we have a beta-cutoff, since while searching this
				// node previously, we found a value greater than the current beta. So we
				// can stop searching and use beta.
				score = beta
				shouldUse = true
			}
		}
	}

	return score, bestMove, shouldUse
}

func (entry *SearchEntry) StoreNewInfo(hash uint64, bestMove uint32, score int16, depth int8, nodeType, ply, age uint8) {
	if score > CheckmateThreshold {
		score += int16(ply)
	} else if score < -CheckmateThreshold {
		score -= int16(ply)
	}

	entry.Hash = hash
	entry.BestMove = bestMove
	entry.Score = score
	entry.Depth = depth

	entry.NodeTypeAndAge = 0
	entry.NodeTypeAndAge |= (nodeType << 6)
	entry.NodeTypeAndAge |= age
}

type PerftEntry struct {
	Hash          uint64
	NodesAndDepth uint64
}

func (entry *PerftEntry) GetNodeCount(hash uint64, depth uint8) (nodeCount uint64, shouldUse bool) {
	if entry.Hash == hash && uint8(entry.NodesAndDepth&RightMost8Bits) == depth {
		return (entry.NodesAndDepth & LeftMost56Bits) >> 8, true
	}
	return 0, false
}

func (entry *PerftEntry) StoreNewInfo(hash, nodes uint64, depth uint8) {
	entry.Hash = hash
	entry.NodesAndDepth = 0
	entry.NodesAndDepth |= (nodes << 8)
	entry.NodesAndDepth |= uint64(depth)
}

type SearchBucket struct {
	entries [BucketSize]SearchEntry
}

func (bucket *SearchBucket) GetEntryForProbing(hash uint64, age uint8) *SearchEntry {
	for i := 0; i < BucketSize; i++ {
		entry := &bucket.entries[i]
		if entry.Hash == hash {
			entry.NodeTypeAndAge &= LeftMost2Bits
			entry.NodeTypeAndAge |= age
			return entry
		}
	}
	return &bucket.entries[BucketSize-1]
}

func (bucket *SearchBucket) GetEntryForStoring(hash uint64, age uint8) *SearchEntry {
	for i := 0; i < BucketSize; i++ {
		if bucket.entries[i].Hash == hash {
			return &bucket.entries[i]
		}
	}

	for i := 0; i < BucketSize; i++ {
		if (bucket.entries[i].NodeTypeAndAge&RightMost6Bits) != age {
			return &bucket.entries[i]
		}
	}

	entryToReplace := &bucket.entries[0]

	for i := 1; i < BucketSize; i++ {
		if bucket.entries[i].Depth < entryToReplace.Depth {
			entryToReplace = &bucket.entries[i]
		}
	}

	return entryToReplace
}

type PerftBucket struct {
	entries [BucketSize]PerftEntry
}

func (bucket *PerftBucket) GetEntryForProbing(hash uint64) *PerftEntry {
	for i := 0; i < BucketSize - 1; i++ {
		if bucket.entries[i].Hash == hash {
			return &bucket.entries[i]
		}
	}
	return &bucket.entries[BucketSize-1]
}

func (bucket *PerftBucket) GetEntryForStoring(hash uint64) *PerftEntry {
	for i := 0; i < BucketSize; i++ {
		if bucket.entries[i].Hash == 0 {
			return &bucket.entries[i]
		}
	}

	entryToReplace := &bucket.entries[0]

	for i := 1; i < BucketSize; i++ {
		entryToReplaceDepth := entryToReplace.NodesAndDepth&RightMost8Bits
		entryDepth := bucket.entries[i].NodesAndDepth&RightMost8Bits
		if entryDepth < entryToReplaceDepth {
			entryToReplace = &bucket.entries[i]
		}
	}

	return entryToReplace
}

type TransTable[Bucket SearchBucket | PerftBucket] struct {
	entries []Bucket
	size    uint64
}

func NewTransTable[Bucket SearchBucket | PerftBucket](sizeInMB uint64) *TransTable[Bucket] {
	tt := TransTable[Bucket]{}
	tt.Resize(sizeInMB)
	return &tt
}

func (tt *TransTable[Bucket]) GetBucket(hash uint64) *Bucket {
	return &tt.entries[hash%tt.size]
}

func (tt *TransTable[Bucket]) Resize(sizeInMB uint64) {
	size := (sizeInMB * 1024 * 1024) / (TTEntrySizeInBytes * BucketSize)
	tt.entries = make([]Bucket, size)
	tt.size = size
}

func (tt *TransTable[Bucket]) Clear() {
	for idx := uint64(0); idx < tt.size; idx++ {
		tt.entries[idx] = *new(Bucket)
	}
}
