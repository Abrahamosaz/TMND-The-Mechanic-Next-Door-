data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres", // | postgres | sqlite | sqlserver
  ]
}


env "gorm" {
  src = data.external_schema.gorm.url
  dev = "postgresql://ab:abraham@localhost:5432/TMND"

  migration {
    dir = "file://cmd/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

