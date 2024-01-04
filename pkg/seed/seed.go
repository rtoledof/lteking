package seed

var seeders = map[string]func() Seeder{}

func RegisterSeeder(name string, fn func() Seeder) {
	seeders[name] = fn
}

type Seeder interface {
	Up() error
	Down() error
}

func Up() error {
	for _, v := range seeders {
		if err := v().Up(); err != nil {
			return err
		}
	}
	return nil
}

func Down() error {
	for _, v := range seeders {
		if err := v().Down(); err != nil {
			return err
		}
	}
	return nil
}
