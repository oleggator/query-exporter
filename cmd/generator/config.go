package main

type Config struct {
	Host     string  `yaml:"host"`
	Port     uint16  `yaml:"port"`
	User     string  `yaml:"user"`
	Password string  `yaml:"password"`
	DBName   string  `yaml:"db_name"`
	Tables   []Table `yaml:"tables"`
}

type Table struct {
	Name   string  `yaml:"name"`
	Count  int     `yaml:"count"`
	Fields []Field `yaml:"fields"`
}

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}
