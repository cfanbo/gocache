package gocache

type CacheOption struct {
	dbcount uint32
	dbindex uint32
	hz      uint32
}

type Option func(*CacheOption)

func NewCacheOption(opts ...Option) *CacheOption {
	defaultOpts := CacheOption{
		dbcount: 16,
		hz:      10,
	}

	// 应用函数参数
	for _, fn := range opts {
		fn(&defaultOpts)
	}

	return &defaultOpts
}

func WithDbNums(num uint32) Option {
	return func(opts *CacheOption) {
		opts.dbcount = num
	}
}

func WithDefaultDb(dbindex uint32) Option {
	return func(opts *CacheOption) {
		opts.dbindex = dbindex
	}
}

func WithHz(hz uint32) Option {
	return func(opts *CacheOption) {
		if hz > 500 {
			hz = 500
		}

		opts.hz = hz
	}
}
