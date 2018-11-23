package main

type Config struct {
	Host      string  `yaml:"host"`
	Port      uint16  `yaml:"port"`
	User      string  `yaml:"user"`
	Password  string  `yaml:"password"`
	DBName    string  `yaml:"db_name"`
	OutputDir string  `yaml:"output_dir"`
	Queries   []Query `yaml:"queries"`
}

type Query struct {
	QueryString   string `yaml:"query"`
	Name          string `yaml:"name"`
	MaxLinesCount int    `yaml:"max_lines"`
}
