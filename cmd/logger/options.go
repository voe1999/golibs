package logger

// Option 是创建logger前的配置。
// 一个Option代表了输出格式、输出位置和输出等级这三项的一种组合。
// 例如，Error及以上级别的日志既要输出到stderr，又要写入err.log，可以创建两个Option。
type Option struct {
	// OutputFormat 日志的输出格式
	OutputFormat OutputFormat
	// WriteTo 日志写到哪。
	WriteTo WriteTo
	// MinLevel 对应日志的最低等级。包含。
	MinLevel Level
	// MaxLevel 对应日志的最高等级。包含。
	MaxLevel Level
}

// OutputFormat 日志的输出格式。
type OutputFormat int

const (
	// FormatPlainText 可读性强的格式。
	FormatPlainText OutputFormat = iota + 1
	// FormatJSON JSON格式。
	FormatJSON
)

// WriteTo 写入位置。
// 目前支持标准io、文件和消息队列。
type WriteTo struct {
	Type WriteToType
	// Path WriteToFile需要额外指定文件路径和文件名。
	// 例如："/var/log/应用名/error.log",
	Path string
}

// WriteToType 日志写入哪里。
type WriteToType int

const (
	// WriteToStdout 写入stdout。
	WriteToStdout WriteToType = iota + 1
	// WriteToStderr 写入stderr
	WriteToStderr
	// WriteToFile 写入日志文件。
	WriteToFile
	// WriteToMessageQueue 写入消息队列。
	WriteToMessageQueue
)

// Level 日志级别。
// 直接采用了zap的级别表示。
type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)
