package engine

const (
	QUEEN_PHASE_SCORE  uint8 = 6
	ROOK_PHASE_SCORE   uint8 = 4
	BISHOP_PHASE_SCORE uint8 = 2
	KNIGHT_PHASE_SCORE uint8 = 2
	PAWN_PHASE_SCORE   uint8 = 0
	TOTAL_PHASE        uint8 = 2*QUEEN_PHASE_SCORE +
							   4*ROOK_PHASE_SCORE +
							   4*BISHOP_PHASE_SCORE +
							   4*KNIGHT_PHASE_SCORE +
							   16*PAWN_PHASE_SCORE
)

var PhaseScores = [6]uint8{
	PAWN_PHASE_SCORE,
	KNIGHT_PHASE_SCORE,
	BISHOP_PHASE_SCORE,
	ROOK_PHASE_SCORE,
	QUEEN_PHASE_SCORE,
}

var MG_PIECE_VALUES = [6]int16{47, 185, 205, 306, 780}
var EG_PIECE_VALUES = [6]int16{132, 289, 306, 493, 808}

var MG_PSQT = [6][64]int16{
    {
          0,   0,   0,   0,   0,   0,   0,   0,
         21,  38,  56,  58,  41,  30,  14,   5,
         -8,  11,  35,  46,  51,  45,  42,  19,
        -29, -16, -16,  -4,   9,   3,  16,  -5,
        -39, -24, -24, -18, -12, -10,   2, -16,
        -44, -29, -35, -36, -20, -13,   6, -17,
        -44, -32, -35, -46, -25, -13,   4, -26,
          0,   0,   0,   0,   0,   0,   0,   0,
    },
    {
        -40,   0,   7,   5,  -2,  -3,  -1,  -4,
         -4,  12,  10,  23,  22,  18,   0,  12,
         -1,  18,  36,  49,  58,  38,  35,   8,
         -3,  -3,   2,  21,   6,  27,   4,  21,
        -21,   3,  -2,  -6,  -1,   4,  13, -12,
        -30, -17, -14,   1,   0,  -7,   2, -16,
        -39, -32, -16, -17, -15, -16, -31, -28,
        -37, -30, -43, -40, -28, -27, -24, -25,
    },
    {
         -8,  -6,   2, -10,  -3,  -6,  -3,  -8,
         -3,  -6,  -3,   0,   5,  10, -15, -10,
         -4,  20,   9,  28,  25,  25,  29,  15,
         -7, -11,   6,  22,   8,  24,  -5,  -4,
         -6,  -3,  -6,   3,   2,  -1,  -1,   0,
        -12,   3,  -5,  -7,  -7, -12,  -2,  -6,
         -5, -12,   4, -18, -16,  -9,  -6, -19,
        -20,  -6, -22, -33, -28, -23, -21, -18,
    },
    {
         15,  15,  12,  18,   7,   6,   6,   9,
          3,  -3,   6,  25,  23,  17,   5,  16,
         -5,  13,   6,  19,  33,  35,  22,   8,
         -2,  -7, -11,  11,   0,  -5,   1,   7,
        -34, -29, -29, -15, -34, -21,  -4, -28,
        -52, -28, -27, -29, -36, -30, -20, -20,
        -50, -38, -28, -29, -35, -37, -22, -41,
        -43, -37, -29, -26, -26, -40, -18, -38,
    },
    {
        -12,  -1,   4,   7,   2,   3,   5,   2,
        -11, -16,   4,  -5,   6,  10,  -3,   6,
         -6,   0,   0,  13,  25,  30,  26,  31,
         -8,   0,   1,  -2,   2,  10,  -2,   7,
        -15, -11, -11, -14, -13,   2,   0,   7,
         -8, -11,  -8, -16,  -3,  -7,   0,  -1,
        -27, -14,  -6,  -3,  -3,   1,  -8,  -9,
         -8, -21, -18, -12,  -8, -33,  -2,   6,
    },
    {
          0,   0,   0,   0,   0,   0,   0,   0,
          1,   3,   2,   2,   1,   1,   4,   0,
          0,   8,   6,   4,   3,   6,   5,   1,
          1,   4,   6,   5,   6,   5,   7,  -1,
         -1,   4,   7,   5,   5,   4,  -1,  -8,
         -7,   4,  11, -11, -12,  -6, -16, -22,
          7,   7,  -7, -31, -24, -24,   5,  -2,
          4,  43,  12, -54,  -1, -35,  21,  22,
    },
}

var EG_PSQT = [6][64]int16{
    {
          0,   0,   0,   0,   0,   0,   0,   0,
         58,  60,  53,  48,  48,  22,  21,  22,
         46,  31,  14,  14,  10,  13,  38,  14,
         14,   4,  -5, -14, -22, -19,  -6,  -8,
         -5,  -3, -28, -35, -26, -23, -18, -25,
        -16,  -5, -27, -17, -19, -25, -18, -31,
         -6,   0, -12,  -7,   6, -11, -14, -29,
          0,   0,   0,   0,   0,   0,   0,   0,
    },
    {
        -21,  -1,   4,   4,  -4,  -2,  -1,  -3,
        -11,   6,   6,  21,  15,   6,  -1,   0,
          1,   7,  25,  20,  27,  13,  18,   6,
         13,  21,  31,  41,  37,  27,  23,  -3,
        -17,   4,  28,  31,  37,  28,   7, -15,
        -36,  -5,   9,  22,  16,  -3, -21, -15,
        -30, -15, -21, -16, -10,  -8, -18, -18,
        -21, -65, -29, -18, -37, -38, -47, -13,
    },
    {
         -4,  -8,   9,  -2,  -1,  -4,  -7,  -4,
         -8,   0,   2,   9,   0,   3,   0, -14,
          4,   8,  10,   2,  17,  15,  13,  18,
         -5,  20,  11,  23,  24,   7,   8,   7,
         -8,  14,  25,  10,  14,  10,   0, -18,
        -14,   9,   9,  17,  16,   4,  -7,  -1,
        -18, -23, -10,   1,   7, -14,  -7, -16,
        -23, -18, -37, -11, -17,  -8,  -8, -17,
    },
    {
          7,  16,  21,  19,  20,   5,  11,  11,
          9,  23,  25,  30,  29,  10,   8,  11,
         21,  15,  19,  19,  14,  13,  18,   1,
          1,   8,  18,  13,   3,   6,  -3,   1,
         -2,  -9,   7,  -3,   1,  -5,  -6, -12,
        -10, -13, -14, -14, -11, -23, -13, -18,
        -17, -21, -15, -11, -26, -26, -29, -13,
        -17, -13,  -6,  -2, -14, -10, -28, -28,
    },
    {
          0,   5,   4,   3,   4,   5,   5,   3,
        -10,   5,  11,   6,   4,   7,   5,  -3,
        -10,   6,  13,  14,  15,  27,  11,  11,
         -6,  -7,  11,  20,  17,  18,  10,   1,
        -11,   0,   3,  25,   4,  -8,   0,  -4,
        -10, -13,   0,   2,  -7,   6,  -9,  -1,
        -12, -12, -23, -19, -32, -33, -15, -10,
        -11,  -8, -29, -28, -29, -26, -10,   1,
    },
    {
          0,   1,   0,   3,   4,   2,   2,  -2,
          1,  12,   9,   6,   5,   8,  14,  -2,
          1,  25,  28,  15,  15,  24,  21,   4,
          7,  19,  29,  28,  34,  22,  23,   2,
          2,  19,  22,  24,  21,  16,   6, -15,
         -7,  -3,   6,   2,   5,  -3, -10, -24,
         -5, -19,  -8,  -7,  -5,  -8, -22, -30,
        -16, -27, -18, -27, -61, -22, -39, -70,
    },
}

var FlipSq = [2][64]uint8{
	{
		0,  1,  2,  3,  4,  5,  6,  7,
		8,  9,  10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
	},

	{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8,  9,  10, 11, 12, 13, 14, 15,
		0,  1,  2,  3,  4,  5,  6,  7,
	},
}

func Evaluate(pos *Position) int16 {
	mgScore := pos.MGScores[pos.SideToMove] - pos.MGScores[pos.SideToMove^1]
	egScore := pos.EGScores[pos.SideToMove] - pos.EGScores[pos.SideToMove^1]

	mgTerm := int32(pos.Phase) * int32(mgScore)
	egTerm := int32(egScore) * int32(TOTAL_PHASE-pos.Phase)
	return int16((mgTerm + egTerm) / int32(TOTAL_PHASE))
}
