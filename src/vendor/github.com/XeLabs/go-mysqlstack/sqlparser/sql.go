//line sql.y:6
package sqlparser

import __yyfmt__ "fmt"

//line sql.y:6
func setParseTree(yylex interface{}, stmt Statement) {
	yylex.(*Tokenizer).ParseTree = stmt
}

func setAllowComments(yylex interface{}, allow bool) {
	yylex.(*Tokenizer).AllowComments = allow
}

func incNesting(yylex interface{}) bool {
	yylex.(*Tokenizer).nesting++
	if yylex.(*Tokenizer).nesting == 200 {
		return true
	}
	return false
}

func decNesting(yylex interface{}) {
	yylex.(*Tokenizer).nesting--
}

func forceEOF(yylex interface{}) {
	yylex.(*Tokenizer).ForceEOF = true
}

//line sql.y:34
type yySymType struct {
	yys         int
	empty       struct{}
	statement   Statement
	selStmt     SelectStatement
	byt         byte
	bytes       []byte
	bytes2      [][]byte
	str         string
	selectExprs SelectExprs
	selectExpr  SelectExpr
	columns     Columns
	colName     *ColName
	tableExprs  TableExprs
	tableExpr   TableExpr
	tableName   *TableName
	indexHints  *IndexHints
	expr        Expr
	boolExpr    BoolExpr
	boolVal     BoolVal
	valExpr     ValExpr
	colTuple    ColTuple
	valExprs    ValExprs
	values      Values
	valTuple    ValTuple
	subquery    *Subquery
	caseExpr    *CaseExpr
	whens       []*When
	when        *When
	orderBy     OrderBy
	order       *Order
	limit       *Limit
	insRows     InsertRows
	updateExprs UpdateExprs
	updateExpr  *UpdateExpr
	colIdent    ColIdent
	colIdents   []ColIdent
	tableIdent  TableIdent
}

const LEX_ERROR = 57346
const UNION = 57347
const SELECT = 57348
const INSERT = 57349
const REPLACE = 57350
const UPDATE = 57351
const DELETE = 57352
const FROM = 57353
const WHERE = 57354
const GROUP = 57355
const HAVING = 57356
const ORDER = 57357
const BY = 57358
const LIMIT = 57359
const OFFSET = 57360
const FOR = 57361
const ALL = 57362
const DISTINCT = 57363
const AS = 57364
const EXISTS = 57365
const ASC = 57366
const DESC = 57367
const INTO = 57368
const DUPLICATE = 57369
const KEY = 57370
const DEFAULT = 57371
const SET = 57372
const LOCK = 57373
const VALUES = 57374
const LAST_INSERT_ID = 57375
const NEXT = 57376
const VALUE = 57377
const JOIN = 57378
const STRAIGHT_JOIN = 57379
const LEFT = 57380
const RIGHT = 57381
const INNER = 57382
const OUTER = 57383
const CROSS = 57384
const NATURAL = 57385
const USE = 57386
const FORCE = 57387
const ON = 57388
const ID = 57389
const HEX = 57390
const STRING = 57391
const INTEGRAL = 57392
const FLOAT = 57393
const HEXNUM = 57394
const VALUE_ARG = 57395
const LIST_ARG = 57396
const COMMENT = 57397
const NULL = 57398
const TRUE = 57399
const FALSE = 57400
const OR = 57401
const AND = 57402
const NOT = 57403
const BETWEEN = 57404
const CASE = 57405
const WHEN = 57406
const THEN = 57407
const ELSE = 57408
const END = 57409
const LE = 57410
const GE = 57411
const NE = 57412
const NULL_SAFE_EQUAL = 57413
const IS = 57414
const LIKE = 57415
const REGEXP = 57416
const IN = 57417
const SHIFT_LEFT = 57418
const SHIFT_RIGHT = 57419
const MOD = 57420
const UNARY = 57421
const INTERVAL = 57422
const JSON_EXTRACT_OP = 57423
const JSON_UNQUOTE_EXTRACT_OP = 57424
const CREATE = 57425
const ALTER = 57426
const DROP = 57427
const RENAME = 57428
const ANALYZE = 57429
const TRUNCATE = 57430
const TABLE = 57431
const INDEX = 57432
const VIEW = 57433
const TO = 57434
const IGNORE = 57435
const IF = 57436
const UNIQUE = 57437
const USING = 57438
const SHOW = 57439
const DESCRIBE = 57440
const EXPLAIN = 57441
const XA = 57442
const PROCESSLIST = 57443
const STATUS = 57444
const QUERYZ = 57445
const TXNZ = 57446
const PARTITION = 57447
const PARTITIONS = 57448
const HASH = 57449
const ENGINE = 57450
const ENGINES = 57451
const DATABASE = 57452
const DATABASES = 57453
const TABLES = 57454
const KILL = 57455
const CURRENT_TIMESTAMP = 57456
const UNUSED = 57457

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LEX_ERROR",
	"UNION",
	"SELECT",
	"INSERT",
	"REPLACE",
	"UPDATE",
	"DELETE",
	"FROM",
	"WHERE",
	"GROUP",
	"HAVING",
	"ORDER",
	"BY",
	"LIMIT",
	"OFFSET",
	"FOR",
	"ALL",
	"DISTINCT",
	"AS",
	"EXISTS",
	"ASC",
	"DESC",
	"INTO",
	"DUPLICATE",
	"KEY",
	"DEFAULT",
	"SET",
	"LOCK",
	"VALUES",
	"LAST_INSERT_ID",
	"NEXT",
	"VALUE",
	"JOIN",
	"STRAIGHT_JOIN",
	"LEFT",
	"RIGHT",
	"INNER",
	"OUTER",
	"CROSS",
	"NATURAL",
	"USE",
	"FORCE",
	"ON",
	"'('",
	"','",
	"')'",
	"ID",
	"HEX",
	"STRING",
	"INTEGRAL",
	"FLOAT",
	"HEXNUM",
	"VALUE_ARG",
	"LIST_ARG",
	"COMMENT",
	"NULL",
	"TRUE",
	"FALSE",
	"OR",
	"AND",
	"NOT",
	"BETWEEN",
	"CASE",
	"WHEN",
	"THEN",
	"ELSE",
	"END",
	"'='",
	"'<'",
	"'>'",
	"LE",
	"GE",
	"NE",
	"NULL_SAFE_EQUAL",
	"IS",
	"LIKE",
	"REGEXP",
	"IN",
	"'|'",
	"'&'",
	"SHIFT_LEFT",
	"SHIFT_RIGHT",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"MOD",
	"'^'",
	"'~'",
	"UNARY",
	"INTERVAL",
	"'.'",
	"JSON_EXTRACT_OP",
	"JSON_UNQUOTE_EXTRACT_OP",
	"CREATE",
	"ALTER",
	"DROP",
	"RENAME",
	"ANALYZE",
	"TRUNCATE",
	"TABLE",
	"INDEX",
	"VIEW",
	"TO",
	"IGNORE",
	"IF",
	"UNIQUE",
	"USING",
	"SHOW",
	"DESCRIBE",
	"EXPLAIN",
	"XA",
	"PROCESSLIST",
	"STATUS",
	"QUERYZ",
	"TXNZ",
	"PARTITION",
	"PARTITIONS",
	"HASH",
	"ENGINE",
	"ENGINES",
	"DATABASE",
	"DATABASES",
	"TABLES",
	"KILL",
	"CURRENT_TIMESTAMP",
	"UNUSED",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 171,
	96, 269,
	-2, 268,
}

const yyNprod = 273
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 1138

var yyAct = [...]int{

	159, 116, 90, 183, 431, 170, 496, 305, 153, 249,
	407, 397, 295, 424, 340, 338, 357, 334, 339, 288,
	62, 248, 3, 257, 350, 203, 188, 178, 179, 91,
	151, 198, 57, 58, 50, 106, 210, 55, 455, 457,
	53, 101, 95, 303, 139, 482, 481, 480, 97, 113,
	92, 98, 61, 59, 60, 51, 342, 56, 239, 240,
	75, 76, 145, 77, 414, 171, 163, 162, 164, 165,
	166, 167, 387, 261, 168, 118, 231, 232, 233, 234,
	228, 176, 123, 247, 228, 107, 108, 109, 110, 112,
	93, 114, 130, 171, 502, 119, 262, 264, 304, 263,
	262, 155, 156, 456, 86, 438, 175, 364, 157, 133,
	158, 264, 144, 99, 154, 264, 436, 104, 105, 335,
	135, 362, 363, 361, 126, 172, 187, 335, 461, 391,
	92, 346, 201, 92, 200, 306, 74, 149, 100, 177,
	207, 174, 128, 263, 262, 173, 131, 145, 263, 262,
	88, 134, 415, 416, 417, 138, 88, 140, 410, 264,
	69, 145, 205, 93, 264, 245, 246, 223, 224, 225,
	147, 360, 148, 88, 94, 469, 470, 382, 70, 66,
	72, 73, 259, 71, 197, 93, 68, 260, 65, 67,
	279, 281, 187, 187, 171, 284, 220, 222, 103, 221,
	92, 293, 291, 200, 351, 353, 354, 137, 297, 352,
	187, 300, 229, 230, 231, 232, 233, 234, 228, 93,
	290, 42, 124, 504, 306, 125, 96, 294, 298, 286,
	331, 306, 301, 227, 226, 235, 236, 229, 230, 231,
	232, 233, 234, 228, 102, 322, 64, 408, 328, 285,
	238, 187, 187, 323, 326, 306, 330, 332, 204, 259,
	187, 344, 328, 395, 306, 306, 348, 349, 187, 187,
	320, 321, 358, 324, 327, 88, 425, 410, 132, 336,
	237, 132, 345, 337, 254, 306, 376, 306, 92, 369,
	371, 280, 355, 22, 132, 374, 368, 378, 373, 258,
	375, 255, 331, 467, 283, 427, 370, 377, 290, 395,
	145, 217, 476, 287, 227, 226, 235, 236, 229, 230,
	231, 232, 233, 234, 228, 215, 383, 93, 199, 260,
	145, 187, 384, 129, 145, 386, 54, 88, 254, 392,
	425, 187, 121, 299, 450, 145, 219, 142, 115, 451,
	344, 390, 448, 385, 479, 478, 380, 449, 447, 381,
	413, 163, 162, 164, 165, 166, 167, 358, 418, 168,
	343, 452, 446, 403, 404, 494, 82, 428, 419, 429,
	22, 359, 84, 196, 195, 426, 433, 495, 435, 81,
	484, 465, 127, 434, 85, 208, 214, 216, 212, 141,
	412, 344, 344, 344, 344, 117, 289, 78, 79, 150,
	194, 443, 184, 445, 442, 435, 444, 453, 193, 458,
	213, 460, 441, 379, 63, 206, 462, 218, 399, 402,
	403, 404, 400, 374, 401, 405, 466, 296, 477, 440,
	111, 394, 204, 89, 474, 501, 187, 241, 242, 243,
	244, 475, 473, 492, 471, 22, 42, 44, 1, 343,
	411, 406, 250, 399, 402, 403, 404, 400, 252, 401,
	405, 256, 485, 209, 52, 486, 359, 302, 211, 282,
	487, 488, 187, 187, 192, 292, 489, 490, 491, 493,
	497, 497, 497, 92, 169, 500, 468, 498, 499, 430,
	160, 439, 393, 507, 503, 508, 505, 506, 509, 389,
	343, 343, 343, 343, 251, 333, 161, 152, 307, 308,
	309, 310, 311, 312, 313, 314, 315, 316, 317, 318,
	319, 372, 122, 265, 185, 454, 398, 184, 184, 396,
	341, 253, 87, 181, 120, 80, 184, 41, 83, 20,
	19, 87, 18, 17, 16, 87, 87, 356, 15, 21,
	365, 366, 367, 145, 14, 13, 171, 163, 162, 164,
	165, 166, 167, 12, 11, 168, 190, 191, 10, 9,
	87, 8, 176, 7, 87, 6, 5, 4, 2, 87,
	0, 136, 0, 87, 0, 87, 0, 0, 143, 0,
	0, 0, 155, 156, 43, 0, 146, 175, 87, 157,
	87, 158, 0, 0, 0, 182, 0, 184, 0, 0,
	0, 0, 87, 388, 0, 202, 172, 0, 45, 46,
	47, 48, 49, 0, 87, 0, 0, 87, 0, 0,
	0, 0, 174, 0, 0, 0, 173, 0, 0, 0,
	250, 0, 0, 463, 420, 421, 422, 0, 0, 0,
	0, 0, 325, 0, 189, 0, 250, 0, 0, 0,
	0, 432, 227, 226, 235, 236, 229, 230, 231, 232,
	233, 234, 228, 437, 93, 0, 0, 87, 145, 0,
	306, 171, 163, 162, 164, 165, 166, 167, 0, 0,
	168, 190, 191, 0, 0, 186, 0, 176, 0, 0,
	0, 0, 0, 0, 0, 464, 227, 226, 235, 236,
	229, 230, 231, 232, 233, 234, 228, 155, 156, 180,
	472, 0, 175, 250, 157, 0, 158, 0, 0, 0,
	182, 182, 329, 0, 0, 0, 189, 0, 0, 182,
	87, 172, 0, 0, 0, 483, 347, 0, 0, 432,
	0, 0, 0, 0, 0, 0, 0, 174, 0, 0,
	145, 173, 306, 171, 163, 162, 164, 165, 166, 167,
	0, 0, 168, 190, 191, 0, 0, 186, 0, 176,
	0, 0, 0, 0, 87, 0, 0, 87, 0, 0,
	0, 0, 423, 0, 0, 0, 0, 0, 0, 155,
	156, 180, 0, 0, 175, 0, 157, 0, 158, 0,
	182, 227, 226, 235, 236, 229, 230, 231, 232, 233,
	234, 228, 189, 172, 0, 0, 0, 409, 0, 87,
	235, 236, 229, 230, 231, 232, 233, 234, 228, 174,
	0, 0, 0, 173, 0, 0, 145, 0, 0, 171,
	163, 162, 164, 165, 166, 167, 0, 22, 168, 190,
	191, 0, 0, 186, 0, 176, 0, 0, 0, 0,
	0, 0, 0, 0, 189, 0, 0, 0, 0, 0,
	87, 87, 87, 87, 0, 155, 156, 180, 0, 0,
	175, 0, 157, 409, 158, 0, 459, 0, 145, 0,
	0, 171, 163, 162, 164, 165, 166, 167, 0, 172,
	168, 190, 191, 0, 0, 186, 0, 176, 0, 0,
	0, 189, 22, 0, 0, 174, 0, 0, 0, 173,
	0, 0, 0, 0, 0, 0, 0, 155, 156, 0,
	0, 0, 175, 0, 157, 145, 158, 0, 171, 163,
	162, 164, 165, 166, 167, 0, 0, 168, 190, 191,
	0, 172, 186, 145, 176, 0, 171, 163, 162, 164,
	165, 166, 167, 0, 0, 168, 0, 174, 0, 0,
	0, 173, 176, 0, 155, 156, 0, 0, 0, 175,
	0, 157, 0, 158, 22, 23, 24, 25, 26, 0,
	0, 0, 155, 156, 0, 0, 0, 175, 172, 157,
	0, 158, 0, 0, 0, 0, 0, 0, 27, 0,
	0, 0, 0, 0, 174, 0, 172, 0, 173, 0,
	0, 0, 36, 0, 0, 0, 0, 0, 0, 0,
	267, 270, 174, 0, 0, 0, 173, 272, 273, 274,
	275, 276, 277, 278, 271, 268, 269, 266, 227, 226,
	235, 236, 229, 230, 231, 232, 233, 234, 228, 227,
	226, 235, 236, 229, 230, 231, 232, 233, 234, 228,
	0, 0, 0, 0, 0, 0, 0, 28, 29, 31,
	30, 33, 32, 0, 0, 0, 0, 0, 0, 0,
	0, 37, 40, 39, 34, 0, 0, 0, 0, 35,
	0, 0, 0, 0, 0, 0, 0, 38, 226, 235,
	236, 229, 230, 231, 232, 233, 234, 228,
}
var yyPact = [...]int{

	998, -1000, -1000, 451, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -71, -72,
	-48, -73, -51, -53, -1000, 408, 196, 61, 83, -1000,
	-1000, 449, 387, 355, -1000, -72, 368, 123, 432, 113,
	-68, -68, -58, -1000, -54, -1000, 123, -69, 194, -69,
	123, 123, -1000, -88, -1000, -1000, -1000, 429, -1000, -56,
	-1000, 302, 388, 388, -1000, -1000, -1000, -1000, -1000, -1000,
	305, 169, -1000, 66, 366, 123, 303, -4, -1000, 123,
	233, -1000, 38, -1000, 123, 56, 123, 157, 123, -64,
	123, 376, 301, 123, -1000, -1000, 263, -1000, -1000, -1000,
	-1000, 123, -1000, 123, -1000, 123, -1000, 15, -1000, -1000,
	809, -1000, 399, -1000, 352, 351, -1000, 123, 298, 113,
	123, 430, 113, 15, 263, 372, -1000, -76, 296, 123,
	-1000, -1000, 123, -1000, 147, -1000, -1000, -1000, -1000, -1000,
	232, -1000, -39, -1000, -1000, 15, 15, 15, 15, 263,
	263, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -13,
	926, -1000, -1000, -1000, -1000, -1000, 15, -1000, 290, -1000,
	-1000, 277, -23, 81, 986, -1000, 908, 861, -1000, 263,
	-1000, -1000, 123, -1000, -1000, -1000, -1000, 283, 374, 113,
	113, 246, -1000, 422, 908, -1000, 997, -1000, -1000, 297,
	113, -1000, -65, 27, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 206, -1000, -1000, -1000, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 310,
	310, -1000, -1000, -1000, 634, 641, 723, 144, 216, 182,
	997, 52, 997, 430, 809, 100, -1000, -1000, 135, -1000,
	-1000, 43, 908, 908, 145, 516, 114, 42, 15, 15,
	15, 145, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 19,
	986, 86, 986, -1000, 449, -1000, 374, 113, -1000, 263,
	451, 233, 238, -1000, 422, 388, 407, 81, -1000, 123,
	-1000, -1000, 123, -1000, 127, -1000, -1000, 756, 1045, -1000,
	-12, -12, -8, -8, -8, -8, 126, 126, 997, 997,
	-1000, -1000, -1000, -1000, 236, 809, -1000, 236, -1000, -24,
	-1000, 15, -1000, 60, -1000, 908, 428, -1000, 261, 427,
	-1000, -1000, 225, 378, 287, -1000, -1000, -32, 19, 33,
	-1000, -1000, 93, -1000, -1000, -1000, 997, -1000, 926, -1000,
	-1000, 114, 15, 15, 15, 997, 997, 739, -1000, -1000,
	294, 230, 257, -1000, 15, -1000, 113, 388, -1000, 15,
	263, -1000, -1000, -1000, -1000, 236, -1000, 113, 997, 46,
	-1000, 15, 37, 425, 406, 100, 100, 100, 100, -1000,
	336, 322, -1000, 316, 308, 335, -6, -1000, 106, -1000,
	-1000, 123, -1000, 215, 40, -1000, -1000, -1000, 182, -1000,
	997, 997, 590, 15, -1000, 364, -1000, 263, -1000, -1000,
	255, -1000, 151, -1000, -1000, -1000, -1000, 997, 15, 422,
	908, 15, 427, 266, 392, -1000, -1000, -1000, -1000, 319,
	-1000, 318, -1000, -1000, -1000, -59, -60, -61, -1000, -1000,
	-1000, -1000, -1000, 15, 997, 362, -1000, 15, -1000, -1000,
	-1000, -1000, 997, 388, 81, 254, 908, 908, -1000, -1000,
	263, 263, 263, 997, 444, -1000, 356, 81, 81, 113,
	113, 113, 113, -1000, 436, 13, 175, -1000, 175, 175,
	233, -1000, 113, -1000, 113, -1000, -1000, 113, -1000, -1000,
}
var yyPgo = [...]int{

	0, 588, 21, 587, 586, 585, 583, 581, 579, 578,
	574, 573, 565, 564, 559, 558, 554, 553, 552, 550,
	549, 604, 548, 547, 545, 544, 27, 28, 543, 541,
	15, 18, 14, 540, 539, 11, 536, 56, 535, 6,
	25, 3, 534, 26, 533, 19, 30, 291, 532, 24,
	16, 9, 531, 8, 114, 517, 516, 515, 17, 514,
	509, 502, 501, 500, 12, 499, 4, 496, 1, 489,
	31, 485, 13, 2, 29, 484, 336, 138, 478, 477,
	474, 473, 0, 23, 471, 494, 10, 461, 460, 20,
	174, 458, 5, 7, 457,
}
var yyR1 = [...]int{

	0, 91, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 2, 2, 2, 3, 3, 4, 4, 5, 6,
	7, 8, 8, 8, 9, 9, 9, 10, 11, 11,
	11, 12, 13, 15, 16, 17, 18, 18, 18, 18,
	18, 18, 18, 18, 18, 18, 19, 20, 14, 94,
	21, 22, 22, 23, 23, 23, 24, 24, 25, 25,
	26, 26, 27, 27, 27, 27, 28, 28, 84, 84,
	84, 83, 83, 29, 29, 30, 30, 31, 31, 32,
	32, 32, 33, 33, 33, 33, 88, 88, 87, 87,
	87, 86, 86, 34, 34, 34, 34, 35, 35, 35,
	35, 36, 36, 37, 37, 38, 38, 38, 38, 39,
	39, 40, 40, 41, 41, 41, 41, 41, 41, 43,
	43, 42, 42, 42, 42, 42, 42, 42, 42, 42,
	42, 42, 42, 42, 49, 49, 49, 49, 49, 49,
	44, 44, 44, 44, 44, 44, 44, 50, 50, 50,
	54, 51, 51, 47, 47, 47, 47, 47, 47, 47,
	47, 47, 47, 47, 47, 47, 47, 47, 47, 47,
	47, 47, 47, 47, 47, 47, 47, 47, 47, 47,
	63, 63, 63, 63, 56, 59, 59, 57, 57, 58,
	60, 60, 55, 55, 55, 46, 46, 46, 46, 46,
	46, 46, 48, 48, 48, 61, 61, 62, 62, 64,
	64, 65, 65, 66, 67, 67, 67, 68, 68, 68,
	68, 69, 69, 69, 70, 70, 71, 71, 72, 72,
	45, 45, 52, 52, 53, 73, 73, 74, 75, 75,
	77, 77, 90, 90, 76, 76, 78, 78, 78, 78,
	78, 78, 79, 79, 80, 80, 81, 81, 82, 85,
	92, 93, 89,
}
var yyR2 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 12, 6, 3, 8, 8, 6, 6, 8, 7,
	3, 6, 4, 9, 6, 7, 7, 5, 4, 5,
	4, 3, 3, 2, 7, 3, 3, 3, 3, 3,
	5, 5, 3, 5, 4, 4, 3, 2, 2, 0,
	2, 0, 2, 1, 2, 2, 0, 1, 0, 1,
	1, 3, 1, 2, 3, 5, 1, 1, 0, 1,
	2, 1, 1, 0, 2, 1, 3, 1, 1, 3,
	3, 3, 3, 5, 5, 3, 0, 1, 0, 1,
	2, 1, 1, 1, 2, 2, 1, 2, 3, 2,
	3, 2, 2, 1, 3, 0, 5, 5, 5, 1,
	3, 0, 2, 1, 3, 3, 2, 3, 3, 1,
	1, 1, 3, 3, 3, 4, 3, 4, 3, 4,
	5, 6, 3, 2, 1, 2, 1, 2, 1, 2,
	1, 1, 1, 1, 1, 1, 1, 3, 1, 1,
	3, 1, 3, 1, 1, 1, 1, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 3, 3, 4, 5, 3, 4, 1,
	1, 1, 1, 1, 5, 0, 1, 1, 2, 4,
	0, 2, 1, 3, 5, 1, 1, 1, 1, 1,
	1, 1, 1, 2, 2, 0, 3, 0, 2, 0,
	3, 1, 3, 2, 0, 1, 1, 0, 2, 4,
	4, 0, 2, 4, 0, 3, 1, 3, 0, 5,
	2, 1, 1, 3, 3, 1, 3, 3, 1, 1,
	0, 2, 0, 3, 0, 1, 1, 1, 1, 1,
	1, 1, 0, 1, 0, 1, 0, 2, 1, 1,
	1, 1, 0,
}
var yyChk = [...]int{

	-1000, -91, -1, -2, -3, -4, -5, -6, -7, -8,
	-9, -10, -11, -12, -13, -15, -16, -17, -18, -19,
	-20, -14, 6, 7, 8, 9, 10, 30, 99, 100,
	102, 101, 104, 103, 116, 121, 44, 113, 129, 115,
	114, -23, 5, -21, -94, -21, -21, -21, -21, -21,
	105, 126, -80, 111, -76, 109, 105, 105, 106, 126,
	105, 105, -89, 16, 50, 127, 118, 128, 125, 99,
	117, 122, 119, 120, 53, -89, -89, -2, 20, 21,
	-24, 34, 21, -22, -76, 26, -37, -85, 50, 11,
	-73, -74, -82, 50, -90, 110, -90, 106, 105, -37,
	-77, 110, 50, -77, -37, -37, 123, -89, -89, -89,
	-89, 11, -89, 105, -89, 46, -68, 17, -68, -89,
	-25, 37, -48, -82, 53, 56, 58, 26, -37, 30,
	96, -37, 48, 71, -37, 64, -85, 50, -37, 108,
	-37, 23, 46, -85, -92, 47, -85, -37, -37, -89,
	-47, -46, -55, -53, -54, 86, 87, 93, 95, -82,
	-63, -56, 52, 51, 53, 54, 55, 56, 59, -85,
	-92, 50, 110, 130, 126, 91, 66, -89, -26, -27,
	88, -28, -85, -41, -47, -42, 64, -92, -43, 23,
	60, 61, -75, 19, 11, 32, 32, -37, -70, 30,
	-92, -73, -85, -40, 12, -74, -47, -92, 23, -81,
	112, -78, 102, 124, 100, 29, 101, 15, 131, 50,
	-37, -37, 50, -89, -89, -89, 83, 82, 92, 86,
	87, 88, 89, 90, 91, 84, 85, 48, 18, 97,
	98, -47, -47, -47, -47, -92, -92, 96, -2, -51,
	-47, -59, -47, -29, 48, 11, -84, -83, 22, -82,
	52, 96, 63, 62, 78, -44, 81, 64, 79, 80,
	65, 78, 71, 72, 73, 74, 75, 76, 77, -41,
	-47, -41, -47, -54, -92, -37, -70, 30, -45, 32,
	-2, -73, -71, -82, -40, -64, 15, -41, -89, 46,
	-82, -89, -79, 108, 71, -93, 49, -47, -47, -47,
	-47, -47, -47, -47, -47, -47, -47, -47, -47, -47,
	-46, -46, -82, -93, -26, 21, -93, -26, -82, -85,
	-93, 48, -93, -57, -58, 67, -40, -27, -30, -31,
	-32, -33, -37, -54, -92, -83, 88, -85, -41, -41,
	-49, 59, 64, 60, 61, -43, -47, -50, -92, -54,
	57, 81, 79, 80, 65, -47, -47, -47, -49, -93,
	-45, -73, -52, -53, -92, -93, 48, -64, -68, 16,
	-37, -37, 50, -89, -93, -26, -93, 96, -47, -60,
	-58, 69, -41, -61, 13, 48, -34, -35, -36, 36,
	40, 42, 37, 38, 39, 43, -87, -86, 22, -85,
	52, -88, 22, -30, 96, 59, 60, 61, -51, -50,
	-47, -47, -47, 63, -72, 46, -72, 48, -82, -68,
	-65, -66, -47, -92, -93, -82, 70, -47, 68, -62,
	14, 16, -31, -32, -31, -32, 36, 36, 36, 41,
	36, 41, 36, -35, -38, 44, 109, 45, -86, -85,
	-93, 88, -93, 63, -47, 27, -53, 48, -67, 24,
	25, -89, -47, -64, -41, -51, 46, 46, 36, 36,
	106, 106, 106, -47, 28, -66, -68, -41, -41, -92,
	-92, -92, 9, -69, 19, 31, -39, -82, -39, -39,
	-73, 9, 81, -93, 48, -93, -93, -82, -82, -82,
}
var yyDef = [...]int{

	0, -2, 1, 2, 3, 4, 5, 6, 7, 8,
	9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
	19, 20, 59, 59, 59, 59, 59, 59, 264, 254,
	0, 0, 0, 0, 272, 0, 0, 0, 0, 272,
	272, 0, 63, 66, 61, 254, 0, 0, 0, 0,
	252, 252, 0, 265, 0, 255, 0, 250, 0, 250,
	0, 0, 43, 0, 272, 272, 272, 272, 272, 0,
	272, 0, 227, 227, 272, 57, 58, 23, 64, 65,
	68, 0, 67, 60, 0, 0, 0, 113, 269, 0,
	30, 245, 0, 268, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 41, 42, 0, 45, 46, 47,
	48, 0, 49, 0, 52, 0, 272, 0, 272, 56,
	0, 69, 0, 212, 0, 0, 62, 0, 234, 0,
	0, 121, 0, 0, 0, 0, 32, 266, 0, 0,
	38, 251, 0, 40, 0, 270, 272, 272, 272, 54,
	228, 163, 164, 165, 166, 0, 0, 0, 0, 202,
	0, 189, 205, 206, 207, 208, 209, 210, 211, 0,
	0, -2, 190, 191, 192, 193, 195, 55, 83, 70,
	72, 78, 0, 76, 77, 123, 0, 0, 131, 0,
	129, 130, 0, 248, 249, 213, 214, 234, 0, 0,
	0, 121, 114, 219, 0, 246, 247, 272, 253, 0,
	0, 272, 262, 0, 256, 257, 258, 259, 260, 261,
	37, 39, 0, 50, 51, 53, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 180, 181, 182, 0, 0, 0, 0, 0, 0,
	161, 0, 196, 121, 0, 0, 73, 79, 0, 81,
	82, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 150, 151, 152, 153, 154, 155, 156, 126,
	0, 0, 161, 143, 0, 22, 0, 0, 26, 0,
	241, 27, 0, 236, 219, 227, 0, 122, 31, 0,
	267, 34, 0, 263, 0, 272, 271, 167, 168, 169,
	170, 171, 172, 173, 174, 175, 176, 177, 229, 230,
	178, 179, 183, 184, 0, 0, 187, 0, 203, 0,
	160, 0, 244, 200, 197, 0, 215, 71, 84, 85,
	87, 88, 98, 96, 0, 80, 74, 0, 124, 125,
	128, 144, 0, 146, 148, 132, 133, 134, 0, 158,
	159, 0, 0, 0, 0, 136, 138, 0, 142, 127,
	238, 238, 240, 242, 0, 235, 0, 227, 29, 0,
	0, 35, 36, 44, 185, 0, 188, 0, 162, 0,
	198, 0, 0, 217, 0, 0, 0, 0, 0, 103,
	0, 0, 106, 0, 0, 0, 115, 99, 0, 101,
	102, 0, 97, 0, 0, 145, 147, 149, 0, 135,
	137, 139, 0, 0, 24, 0, 25, 0, 237, 28,
	220, 221, 224, 272, 186, 204, 194, 201, 0, 219,
	0, 0, 86, 92, 0, 95, 104, 105, 107, 0,
	109, 0, 111, 112, 89, 0, 0, 0, 100, 90,
	91, 75, 157, 0, 140, 0, 243, 0, 223, 225,
	226, 33, 199, 227, 218, 216, 0, 0, 108, 110,
	0, 0, 0, 141, 0, 222, 231, 93, 94, 0,
	0, 0, 0, 21, 0, 0, 0, 119, 0, 0,
	239, 232, 0, 116, 0, 117, 118, 0, 120, 233,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 90, 83, 3,
	47, 49, 88, 86, 48, 87, 96, 89, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	72, 71, 73, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 92, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 82, 3, 93,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 50, 51, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 61, 62, 63, 64,
	65, 66, 67, 68, 69, 70, 74, 75, 76, 77,
	78, 79, 80, 81, 84, 85, 91, 94, 95, 97,
	98, 99, 100, 101, 102, 103, 104, 105, 106, 107,
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117,
	118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
	128, 129, 130, 131,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:182
		{
			setParseTree(yylex, yyDollar[1].statement)
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:188
		{
			yyVAL.statement = yyDollar[1].selStmt
		}
	case 21:
		yyDollar = yyS[yypt-12 : yypt+1]
		//line sql.y:212
		{
			yyVAL.selStmt = &Select{Comments: Comments(yyDollar[2].bytes2), Distinct: yyDollar[3].str, Hints: yyDollar[4].str, SelectExprs: yyDollar[5].selectExprs, From: yyDollar[6].tableExprs, Where: NewWhere(WhereStr, yyDollar[7].boolExpr), GroupBy: GroupBy(yyDollar[8].valExprs), Having: NewWhere(HavingStr, yyDollar[9].boolExpr), OrderBy: yyDollar[10].orderBy, Limit: yyDollar[11].limit, Lock: yyDollar[12].str}
		}
	case 22:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:216
		{
			yyVAL.selStmt = &Select{Comments: Comments(yyDollar[2].bytes2), SelectExprs: SelectExprs{Nextval{Expr: yyDollar[4].valExpr}}, From: TableExprs{&AliasedTableExpr{Expr: yyDollar[6].tableName}}}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:220
		{
			yyVAL.selStmt = &Union{Type: yyDollar[2].str, Left: yyDollar[1].selStmt, Right: yyDollar[3].selStmt}
		}
	case 24:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line sql.y:226
		{
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Ignore: yyDollar[3].str, Table: yyDollar[5].tableName, Columns: yyDollar[6].columns, Rows: yyDollar[7].insRows, OnDup: OnDup(yyDollar[8].updateExprs)}
		}
	case 25:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line sql.y:230
		{
			cols := make(Columns, 0, len(yyDollar[7].updateExprs))
			vals := make(ValTuple, 0, len(yyDollar[7].updateExprs))
			for _, updateList := range yyDollar[7].updateExprs {
				cols = append(cols, updateList.Name)
				vals = append(vals, updateList.Expr)
			}
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Ignore: yyDollar[3].str, Table: yyDollar[5].tableName, Columns: cols, Rows: Values{vals}, OnDup: OnDup(yyDollar[8].updateExprs)}
		}
	case 26:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:242
		{
			yyVAL.statement = &Replace{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: yyDollar[5].columns, Rows: yyDollar[6].insRows}
		}
	case 27:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:246
		{
			cols := make(Columns, 0, len(yyDollar[6].updateExprs))
			vals := make(ValTuple, 0, len(yyDollar[6].updateExprs))
			for _, updateList := range yyDollar[6].updateExprs {
				cols = append(cols, updateList.Name)
				vals = append(vals, updateList.Expr)
			}
			yyVAL.statement = &Replace{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: cols, Rows: Values{vals}}
		}
	case 28:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line sql.y:258
		{
			yyVAL.statement = &Update{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[3].tableName, Exprs: yyDollar[5].updateExprs, Where: NewWhere(WhereStr, yyDollar[6].boolExpr), OrderBy: yyDollar[7].orderBy, Limit: yyDollar[8].limit}
		}
	case 29:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:264
		{
			yyVAL.statement = &Delete{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Where: NewWhere(WhereStr, yyDollar[5].boolExpr), OrderBy: yyDollar[6].orderBy, Limit: yyDollar[7].limit}
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:270
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: yyDollar[3].updateExprs}
		}
	case 31:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:276
		{
			yyVAL.statement = &DDL{Action: CreateTableStr, IfNotExists: bool(yyDollar[3].boolVal), Table: yyDollar[4].tableName, NewName: yyDollar[4].tableName}
		}
	case 32:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:281
		{
			yyVAL.statement = &DDL{Action: CreateDBStr, IfNotExists: bool(yyDollar[3].boolVal), Database: yyDollar[4].tableIdent}
		}
	case 33:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line sql.y:285
		{
			yyVAL.statement = &DDL{Action: CreateIndexStr, Table: yyDollar[7].tableName, NewName: yyDollar[7].tableName}
		}
	case 34:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:291
		{
			yyVAL.statement = &DDL{Action: AlterStr, Table: yyDollar[4].tableName, NewName: yyDollar[4].tableName}
		}
	case 35:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:295
		{
			yyVAL.statement = &DDL{Action: AlterStr, Table: yyDollar[4].tableName, NewName: yyDollar[7].tableName}
		}
	case 36:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:299
		{
			yyVAL.statement = &DDL{Action: AlterEngineStr, Table: yyDollar[4].tableName, Engine: string(yyDollar[7].bytes)}
		}
	case 37:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:305
		{
			yyVAL.statement = &DDL{Action: RenameStr, Table: yyDollar[3].tableName, NewName: yyDollar[5].tableName}
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:311
		{
			var exists bool
			if yyDollar[3].byt != 0 {
				exists = true
			}
			yyVAL.statement = &DDL{Action: DropTableStr, Table: yyDollar[4].tableName, IfExists: exists}
		}
	case 39:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:319
		{
			yyVAL.statement = &DDL{Action: DropIndexStr, Table: yyDollar[5].tableName, NewName: yyDollar[5].tableName}
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:323
		{
			var exists bool
			if yyDollar[3].byt != 0 {
				exists = true
			}
			yyVAL.statement = &DDL{Action: DropDBStr, Database: yyDollar[4].tableIdent, IfExists: exists}
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:333
		{
			yyVAL.statement = &DDL{Action: TruncateTableStr, Table: yyDollar[3].tableName}
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:339
		{
			yyVAL.statement = &DDL{Action: AlterStr, Table: yyDollar[3].tableName, NewName: yyDollar[3].tableName}
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:345
		{
			yyVAL.statement = &Xa{}
		}
	case 44:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:351
		{
			yyVAL.statement = &Shard{ShardKey: string(yyDollar[5].bytes)}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:357
		{
			yyVAL.statement = &UseDB{Database: string(yyDollar[2].bytes)}
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:363
		{
			yyVAL.statement = &ShowDatabases{}
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:367
		{
			yyVAL.statement = &ShowStatus{}
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:371
		{
			yyVAL.statement = &ShowTables{}
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:375
		{
			yyVAL.statement = &ShowEngines{}
		}
	case 50:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:379
		{
			yyVAL.statement = &ShowTables{Database: yyDollar[4].tableIdent}
		}
	case 51:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:383
		{
			yyVAL.statement = &ShowCreateTable{Table: yyDollar[4].tableName}
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:387
		{
			yyVAL.statement = &ShowProcesslist{}
		}
	case 53:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:391
		{
			yyVAL.statement = &ShowPartitions{Table: yyDollar[4].tableName}
		}
	case 54:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:395
		{
			yyVAL.statement = &ShowQueryz{Limit: yyDollar[3].limit}
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:399
		{
			yyVAL.statement = &ShowTxnz{Limit: yyDollar[3].limit}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:405
		{
			yyVAL.statement = &Kill{QueryID: string(yyDollar[2].bytes)}
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:411
		{
			yyVAL.statement = &Explain{}
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:417
		{
			yyVAL.statement = &Other{}
		}
	case 59:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:422
		{
			setAllowComments(yylex, true)
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:426
		{
			yyVAL.bytes2 = yyDollar[2].bytes2
			setAllowComments(yylex, false)
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:432
		{
			yyVAL.bytes2 = nil
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:436
		{
			yyVAL.bytes2 = append(yyDollar[1].bytes2, yyDollar[2].bytes)
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:442
		{
			yyVAL.str = UnionStr
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:446
		{
			yyVAL.str = UnionAllStr
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:450
		{
			yyVAL.str = UnionDistinctStr
		}
	case 66:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:455
		{
			yyVAL.str = ""
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:459
		{
			yyVAL.str = DistinctStr
		}
	case 68:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:464
		{
			yyVAL.str = ""
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:468
		{
			yyVAL.str = StraightJoinHint
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:474
		{
			yyVAL.selectExprs = SelectExprs{yyDollar[1].selectExpr}
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:478
		{
			yyVAL.selectExprs = append(yyVAL.selectExprs, yyDollar[3].selectExpr)
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:484
		{
			yyVAL.selectExpr = &StarExpr{}
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:488
		{
			yyVAL.selectExpr = &NonStarExpr{Expr: yyDollar[1].expr, As: yyDollar[2].colIdent}
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:492
		{
			yyVAL.selectExpr = &StarExpr{TableName: &TableName{Name: yyDollar[1].tableIdent}}
		}
	case 75:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:496
		{
			yyVAL.selectExpr = &StarExpr{TableName: &TableName{Qualifier: yyDollar[1].tableIdent, Name: yyDollar[3].tableIdent}}
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:502
		{
			yyVAL.expr = yyDollar[1].boolExpr
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:506
		{
			yyVAL.expr = yyDollar[1].valExpr
		}
	case 78:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:511
		{
			yyVAL.colIdent = ColIdent{}
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:515
		{
			yyVAL.colIdent = yyDollar[1].colIdent
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:519
		{
			yyVAL.colIdent = yyDollar[2].colIdent
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:526
		{
			yyVAL.colIdent = NewColIdent(string(yyDollar[1].bytes))
		}
	case 83:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:531
		{
			yyVAL.tableExprs = TableExprs{&AliasedTableExpr{Expr: &TableName{Name: NewTableIdent("dual")}}}
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:535
		{
			yyVAL.tableExprs = yyDollar[2].tableExprs
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:541
		{
			yyVAL.tableExprs = TableExprs{yyDollar[1].tableExpr}
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:545
		{
			yyVAL.tableExprs = append(yyVAL.tableExprs, yyDollar[3].tableExpr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:555
		{
			yyVAL.tableExpr = &AliasedTableExpr{Expr: yyDollar[1].tableName, As: yyDollar[2].tableIdent, Hints: yyDollar[3].indexHints}
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:559
		{
			yyVAL.tableExpr = &AliasedTableExpr{Expr: yyDollar[1].subquery, As: yyDollar[3].tableIdent}
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:563
		{
			yyVAL.tableExpr = &ParenTableExpr{Exprs: yyDollar[2].tableExprs}
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:576
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr}
		}
	case 93:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:580
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr, On: yyDollar[5].boolExpr}
		}
	case 94:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:584
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr, On: yyDollar[5].boolExpr}
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:588
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr}
		}
	case 96:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:593
		{
			yyVAL.empty = struct{}{}
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:595
		{
			yyVAL.empty = struct{}{}
		}
	case 98:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:598
		{
			yyVAL.tableIdent = NewTableIdent("")
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:602
		{
			yyVAL.tableIdent = yyDollar[1].tableIdent
		}
	case 100:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:606
		{
			yyVAL.tableIdent = yyDollar[2].tableIdent
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:613
		{
			yyVAL.tableIdent = NewTableIdent(string(yyDollar[1].bytes))
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:619
		{
			yyVAL.str = JoinStr
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:623
		{
			yyVAL.str = JoinStr
		}
	case 105:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:627
		{
			yyVAL.str = JoinStr
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:631
		{
			yyVAL.str = StraightJoinStr
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:637
		{
			yyVAL.str = LeftJoinStr
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:641
		{
			yyVAL.str = LeftJoinStr
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:645
		{
			yyVAL.str = RightJoinStr
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:649
		{
			yyVAL.str = RightJoinStr
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:655
		{
			yyVAL.str = NaturalJoinStr
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:659
		{
			if yyDollar[2].str == LeftJoinStr {
				yyVAL.str = NaturalLeftJoinStr
			} else {
				yyVAL.str = NaturalRightJoinStr
			}
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:669
		{
			yyVAL.tableName = &TableName{Name: yyDollar[1].tableIdent}
		}
	case 114:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:673
		{
			yyVAL.tableName = &TableName{Qualifier: yyDollar[1].tableIdent, Name: yyDollar[3].tableIdent}
		}
	case 115:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:678
		{
			yyVAL.indexHints = nil
		}
	case 116:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:682
		{
			yyVAL.indexHints = &IndexHints{Type: UseStr, Indexes: yyDollar[4].colIdents}
		}
	case 117:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:686
		{
			yyVAL.indexHints = &IndexHints{Type: IgnoreStr, Indexes: yyDollar[4].colIdents}
		}
	case 118:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:690
		{
			yyVAL.indexHints = &IndexHints{Type: ForceStr, Indexes: yyDollar[4].colIdents}
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:696
		{
			yyVAL.colIdents = []ColIdent{yyDollar[1].colIdent}
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:700
		{
			yyVAL.colIdents = append(yyDollar[1].colIdents, yyDollar[3].colIdent)
		}
	case 121:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:705
		{
			yyVAL.boolExpr = nil
		}
	case 122:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:709
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 124:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:716
		{
			yyVAL.boolExpr = &AndExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:720
		{
			yyVAL.boolExpr = &OrExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 126:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:724
		{
			yyVAL.boolExpr = &NotExpr{Expr: yyDollar[2].boolExpr}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:728
		{
			yyVAL.boolExpr = &ParenBoolExpr{Expr: yyDollar[2].boolExpr}
		}
	case 128:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:732
		{
			yyVAL.boolExpr = &IsExpr{Operator: yyDollar[3].str, Expr: yyDollar[1].boolExpr}
		}
	case 129:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:738
		{
			yyVAL.boolVal = BoolVal(true)
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:742
		{
			yyVAL.boolVal = BoolVal(false)
		}
	case 131:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:748
		{
			yyVAL.boolExpr = yyDollar[1].boolVal
		}
	case 132:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:752
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: yyDollar[2].str, Right: yyDollar[3].boolVal}
		}
	case 133:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:756
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: yyDollar[2].str, Right: yyDollar[3].valExpr}
		}
	case 134:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:760
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: InStr, Right: yyDollar[3].colTuple}
		}
	case 135:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:764
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: NotInStr, Right: yyDollar[4].colTuple}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:768
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: LikeStr, Right: yyDollar[3].valExpr}
		}
	case 137:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:772
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: NotLikeStr, Right: yyDollar[4].valExpr}
		}
	case 138:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:776
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: RegexpStr, Right: yyDollar[3].valExpr}
		}
	case 139:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:780
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: NotRegexpStr, Right: yyDollar[4].valExpr}
		}
	case 140:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:784
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: BetweenStr, From: yyDollar[3].valExpr, To: yyDollar[5].valExpr}
		}
	case 141:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:788
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: NotBetweenStr, From: yyDollar[4].valExpr, To: yyDollar[6].valExpr}
		}
	case 142:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:792
		{
			yyVAL.boolExpr = &IsExpr{Operator: yyDollar[3].str, Expr: yyDollar[1].valExpr}
		}
	case 143:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:796
		{
			yyVAL.boolExpr = &ExistsExpr{Subquery: yyDollar[2].subquery}
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:802
		{
			yyVAL.str = IsNullStr
		}
	case 145:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:806
		{
			yyVAL.str = IsNotNullStr
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:810
		{
			yyVAL.str = IsTrueStr
		}
	case 147:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:814
		{
			yyVAL.str = IsNotTrueStr
		}
	case 148:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:818
		{
			yyVAL.str = IsFalseStr
		}
	case 149:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:822
		{
			yyVAL.str = IsNotFalseStr
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:828
		{
			yyVAL.str = EqualStr
		}
	case 151:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:832
		{
			yyVAL.str = LessThanStr
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:836
		{
			yyVAL.str = GreaterThanStr
		}
	case 153:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:840
		{
			yyVAL.str = LessEqualStr
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:844
		{
			yyVAL.str = GreaterEqualStr
		}
	case 155:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:848
		{
			yyVAL.str = NotEqualStr
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:852
		{
			yyVAL.str = NullSafeEqualStr
		}
	case 157:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:858
		{
			yyVAL.colTuple = ValTuple(yyDollar[2].valExprs)
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:862
		{
			yyVAL.colTuple = yyDollar[1].subquery
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:866
		{
			yyVAL.colTuple = ListArg(yyDollar[1].bytes)
		}
	case 160:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:872
		{
			yyVAL.subquery = &Subquery{yyDollar[2].selStmt}
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:878
		{
			yyVAL.valExprs = ValExprs{yyDollar[1].valExpr}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:882
		{
			yyVAL.valExprs = append(yyDollar[1].valExprs, yyDollar[3].valExpr)
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:888
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 164:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:892
		{
			yyVAL.valExpr = yyDollar[1].colName
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:896
		{
			yyVAL.valExpr = yyDollar[1].valTuple
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:900
		{
			yyVAL.valExpr = yyDollar[1].subquery
		}
	case 167:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:904
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: BitAndStr, Right: yyDollar[3].valExpr}
		}
	case 168:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:908
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: BitOrStr, Right: yyDollar[3].valExpr}
		}
	case 169:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:912
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: BitXorStr, Right: yyDollar[3].valExpr}
		}
	case 170:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:916
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: PlusStr, Right: yyDollar[3].valExpr}
		}
	case 171:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:920
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: MinusStr, Right: yyDollar[3].valExpr}
		}
	case 172:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:924
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: MultStr, Right: yyDollar[3].valExpr}
		}
	case 173:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:928
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: DivStr, Right: yyDollar[3].valExpr}
		}
	case 174:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:932
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: ModStr, Right: yyDollar[3].valExpr}
		}
	case 175:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:936
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: ModStr, Right: yyDollar[3].valExpr}
		}
	case 176:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:940
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: ShiftLeftStr, Right: yyDollar[3].valExpr}
		}
	case 177:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:944
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: ShiftRightStr, Right: yyDollar[3].valExpr}
		}
	case 178:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:948
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].colName, Operator: JSONExtractOp, Right: yyDollar[3].valExpr}
		}
	case 179:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:952
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].colName, Operator: JSONUnquoteExtractOp, Right: yyDollar[3].valExpr}
		}
	case 180:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:956
		{
			if num, ok := yyDollar[2].valExpr.(*SQLVal); ok && num.Type == IntVal {
				yyVAL.valExpr = num
			} else {
				yyVAL.valExpr = &UnaryExpr{Operator: UPlusStr, Expr: yyDollar[2].valExpr}
			}
		}
	case 181:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:964
		{
			if num, ok := yyDollar[2].valExpr.(*SQLVal); ok && num.Type == IntVal {
				// Handle double negative
				if num.Val[0] == '-' {
					num.Val = num.Val[1:]
					yyVAL.valExpr = num
				} else {
					yyVAL.valExpr = NewIntVal(append([]byte("-"), num.Val...))
				}
			} else {
				yyVAL.valExpr = &UnaryExpr{Operator: UMinusStr, Expr: yyDollar[2].valExpr}
			}
		}
	case 182:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:978
		{
			yyVAL.valExpr = &UnaryExpr{Operator: TildaStr, Expr: yyDollar[2].valExpr}
		}
	case 183:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:982
		{
			// This rule prevents the usage of INTERVAL
			// as a function. If support is needed for that,
			// we'll need to revisit this. The solution
			// will be non-trivial because of grammar conflicts.
			yyVAL.valExpr = &IntervalExpr{Expr: yyDollar[2].valExpr, Unit: yyDollar[3].colIdent}
		}
	case 184:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:990
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].colIdent}
		}
	case 185:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:994
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].colIdent, Exprs: yyDollar[3].selectExprs}
		}
	case 186:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:998
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].colIdent, Distinct: true, Exprs: yyDollar[4].selectExprs}
		}
	case 187:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1002
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].colIdent}
		}
	case 188:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:1006
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].colIdent, Exprs: yyDollar[3].selectExprs}
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1010
		{
			yyVAL.valExpr = yyDollar[1].caseExpr
		}
	case 190:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1016
		{
			yyVAL.colIdent = NewColIdent("if")
		}
	case 191:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1020
		{
			yyVAL.colIdent = NewColIdent("current_timestamp")
		}
	case 192:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1024
		{
			yyVAL.colIdent = NewColIdent("database")
		}
	case 193:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1028
		{
			yyVAL.colIdent = NewColIdent("mod")
		}
	case 194:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:1034
		{
			yyVAL.caseExpr = &CaseExpr{Expr: yyDollar[2].valExpr, Whens: yyDollar[3].whens, Else: yyDollar[4].valExpr}
		}
	case 195:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1039
		{
			yyVAL.valExpr = nil
		}
	case 196:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1043
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 197:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1049
		{
			yyVAL.whens = []*When{yyDollar[1].when}
		}
	case 198:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1053
		{
			yyVAL.whens = append(yyDollar[1].whens, yyDollar[2].when)
		}
	case 199:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:1059
		{
			yyVAL.when = &When{Cond: yyDollar[2].boolExpr, Val: yyDollar[4].valExpr}
		}
	case 200:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1064
		{
			yyVAL.valExpr = nil
		}
	case 201:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1068
		{
			yyVAL.valExpr = yyDollar[2].valExpr
		}
	case 202:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1074
		{
			yyVAL.colName = &ColName{Name: yyDollar[1].colIdent}
		}
	case 203:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1078
		{
			yyVAL.colName = &ColName{Qualifier: &TableName{Name: yyDollar[1].tableIdent}, Name: yyDollar[3].colIdent}
		}
	case 204:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:1082
		{
			yyVAL.colName = &ColName{Qualifier: &TableName{Qualifier: yyDollar[1].tableIdent, Name: yyDollar[3].tableIdent}, Name: yyDollar[5].colIdent}
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1088
		{
			yyVAL.valExpr = NewStrVal(yyDollar[1].bytes)
		}
	case 206:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1092
		{
			yyVAL.valExpr = NewHexVal(yyDollar[1].bytes)
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1096
		{
			yyVAL.valExpr = NewIntVal(yyDollar[1].bytes)
		}
	case 208:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1100
		{
			yyVAL.valExpr = NewFloatVal(yyDollar[1].bytes)
		}
	case 209:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1104
		{
			yyVAL.valExpr = NewHexNum(yyDollar[1].bytes)
		}
	case 210:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1108
		{
			yyVAL.valExpr = NewValArg(yyDollar[1].bytes)
		}
	case 211:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1112
		{
			yyVAL.valExpr = &NullVal{}
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1118
		{
			// TODO(sougou): Deprecate this construct.
			if yyDollar[1].colIdent.Lowered() != "value" {
				yylex.Error("expecting value after next")
				return 1
			}
			yyVAL.valExpr = NewIntVal([]byte("1"))
		}
	case 213:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1127
		{
			yyVAL.valExpr = NewIntVal(yyDollar[1].bytes)
		}
	case 214:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1131
		{
			yyVAL.valExpr = NewValArg(yyDollar[1].bytes)
		}
	case 215:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1136
		{
			yyVAL.valExprs = nil
		}
	case 216:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1140
		{
			yyVAL.valExprs = yyDollar[3].valExprs
		}
	case 217:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1145
		{
			yyVAL.boolExpr = nil
		}
	case 218:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1149
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 219:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1154
		{
			yyVAL.orderBy = nil
		}
	case 220:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1158
		{
			yyVAL.orderBy = yyDollar[3].orderBy
		}
	case 221:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1164
		{
			yyVAL.orderBy = OrderBy{yyDollar[1].order}
		}
	case 222:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1168
		{
			yyVAL.orderBy = append(yyDollar[1].orderBy, yyDollar[3].order)
		}
	case 223:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1174
		{
			yyVAL.order = &Order{Expr: yyDollar[1].valExpr, Direction: yyDollar[2].str}
		}
	case 224:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1179
		{
			yyVAL.str = AscScr
		}
	case 225:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1183
		{
			yyVAL.str = AscScr
		}
	case 226:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1187
		{
			yyVAL.str = DescScr
		}
	case 227:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1192
		{
			yyVAL.limit = nil
		}
	case 228:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1196
		{
			yyVAL.limit = &Limit{Rowcount: yyDollar[2].valExpr}
		}
	case 229:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:1200
		{
			yyVAL.limit = &Limit{Offset: yyDollar[2].valExpr, Rowcount: yyDollar[4].valExpr}
		}
	case 230:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:1204
		{
			yyVAL.limit = &Limit{Rowcount: yyDollar[2].valExpr, Offset: yyDollar[4].valExpr}
		}
	case 231:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1210
		{
			yyVAL.str = ""
		}
	case 232:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1214
		{
			yyVAL.str = ForUpdateStr
		}
	case 233:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:1218
		{
			if yyDollar[3].colIdent.Lowered() != "share" {
				yylex.Error("expecting share")
				return 1
			}
			if yyDollar[4].colIdent.Lowered() != "mode" {
				yylex.Error("expecting mode")
				return 1
			}
			yyVAL.str = ShareModeStr
		}
	case 234:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1231
		{
			yyVAL.columns = nil
		}
	case 235:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1235
		{
			yyVAL.columns = yyDollar[2].columns
		}
	case 236:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1241
		{
			yyVAL.columns = Columns{yyDollar[1].colIdent}
		}
	case 237:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1245
		{
			yyVAL.columns = append(yyVAL.columns, yyDollar[3].colIdent)
		}
	case 238:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1250
		{
			yyVAL.updateExprs = nil
		}
	case 239:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:1254
		{
			yyVAL.updateExprs = yyDollar[5].updateExprs
		}
	case 240:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1260
		{
			yyVAL.insRows = yyDollar[2].values
		}
	case 241:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1264
		{
			yyVAL.insRows = yyDollar[1].selStmt
		}
	case 242:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1270
		{
			yyVAL.values = Values{yyDollar[1].valTuple}
		}
	case 243:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1274
		{
			yyVAL.values = append(yyDollar[1].values, yyDollar[3].valTuple)
		}
	case 244:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1280
		{
			yyVAL.valTuple = ValTuple(yyDollar[2].valExprs)
		}
	case 245:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1286
		{
			yyVAL.updateExprs = UpdateExprs{yyDollar[1].updateExpr}
		}
	case 246:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1290
		{
			yyVAL.updateExprs = append(yyDollar[1].updateExprs, yyDollar[3].updateExpr)
		}
	case 247:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1296
		{
			yyVAL.updateExpr = &UpdateExpr{Name: yyDollar[1].colIdent, Expr: yyDollar[3].valExpr}
		}
	case 250:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1305
		{
			yyVAL.byt = 0
		}
	case 251:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1307
		{
			yyVAL.byt = 1
		}
	case 252:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1310
		{
			yyVAL.boolVal = false
		}
	case 253:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1312
		{
			yyVAL.boolVal = true
		}
	case 254:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1315
		{
			yyVAL.str = ""
		}
	case 255:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1317
		{
			yyVAL.str = IgnoreStr
		}
	case 256:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1321
		{
			yyVAL.empty = struct{}{}
		}
	case 257:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1323
		{
			yyVAL.empty = struct{}{}
		}
	case 258:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1325
		{
			yyVAL.empty = struct{}{}
		}
	case 259:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1327
		{
			yyVAL.empty = struct{}{}
		}
	case 260:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1329
		{
			yyVAL.empty = struct{}{}
		}
	case 261:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1331
		{
			yyVAL.empty = struct{}{}
		}
	case 262:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1334
		{
			yyVAL.empty = struct{}{}
		}
	case 263:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1336
		{
			yyVAL.empty = struct{}{}
		}
	case 264:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1339
		{
			yyVAL.empty = struct{}{}
		}
	case 265:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1341
		{
			yyVAL.empty = struct{}{}
		}
	case 266:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1344
		{
			yyVAL.empty = struct{}{}
		}
	case 267:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1346
		{
			yyVAL.empty = struct{}{}
		}
	case 268:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1350
		{
			yyVAL.colIdent = NewColIdent(string(yyDollar[1].bytes))
		}
	case 269:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1356
		{
			yyVAL.tableIdent = NewTableIdent(string(yyDollar[1].bytes))
		}
	case 270:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1362
		{
			if incNesting(yylex) {
				yylex.Error("max nesting level reached")
				return 1
			}
		}
	case 271:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1371
		{
			decNesting(yylex)
		}
	case 272:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1376
		{
			forceEOF(yylex)
		}
	}
	goto yystack /* stack new state and value */
}
