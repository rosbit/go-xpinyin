package xpinyin

const (
	defSplitter      = "-"
	defNCombinations = 10
)

type options struct {
	toneMarks bool
	toneNumbers bool
	splitter  string
	retroflex bool
	toUpper   bool
	toCapitalize bool
	maxnCombinations int

	splitterSet bool
}

type Option func(*options)

func WithToneMarks() Option {
	return func(o *options) {
		o.toneMarks = true
	}
}

func WithToneNumbers() Option {
	return func(o *options) {
		o.toneNumbers = true
	}
}

func WithSplitter(splitter string) Option {
	return func(o *options) {
		o.splitter = splitter
		o.splitterSet = true
	}
}

func WithRetroflex() Option {
	return func(o *options) {
		o.retroflex = true
	}
}

func ToUpper() Option {
	return func(o *options) {
		o.toUpper = true
	}
}

func ToCapitalize() Option {
	return func(o *options) {
		o.toCapitalize = true
	}
}

func MaxNCombinations(n int) Option {
	return func(o *options) {
		o.maxnCombinations = n
	}
}

func getOptions(opts ...Option) *options {
	var os options
	for _, opt := range opts {
		opt(&os)
	}
	if !os.splitterSet && len(os.splitter) == 0 {
		os.splitter = defSplitter
	}
	if os.maxnCombinations <= 0 {
		os.maxnCombinations = defNCombinations
	}

	return &os
}
