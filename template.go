package mathxf

type template struct {
	tpl  string
	opts *templateOptions

	// Output
	root *nodeDocument
}
type templateOptions struct {
	name string // template name

	bannedTags map[string]struct{}
	debug      bool
}
type Option func(*templateOptions)

func Debug(b bool) Option {
	return func(opts *templateOptions) {
		opts.debug = b
	}
}
func Name(name string) Option {
	return func(opts *templateOptions) {
		opts.name = name
	}
}

func BannedTags(tags ...string) Option {
	return func(opts *templateOptions) {
		for _, tag := range tags {
			opts.bannedTags[tag] = struct{}{}
		}
	}
}

func NewTemplate(tpl string, opts ...Option) *template {
	t := &template{
		tpl: tpl,
		opts: &templateOptions{
			bannedTags: make(map[string]struct{}),
		},
	}
	for _, opt := range opts {
		opt(t.opts)
	}

	return t
}

func (t *template) Name() string {
	return t.opts.name
}

func (t *template) Execute(data interface{}) (ResValues, error) {

	return nil, nil
}
