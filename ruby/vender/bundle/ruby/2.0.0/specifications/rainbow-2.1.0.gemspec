# -*- encoding: utf-8 -*-

Gem::Specification.new do |s|
  s.name = "rainbow"
  s.version = "2.1.0"

  s.required_rubygems_version = Gem::Requirement.new(">= 0") if s.respond_to? :required_rubygems_version=
  s.authors = ["Marcin Kulik"]
  s.date = "2016-01-24"
  s.description = "Colorize printed text on ANSI terminals"
  s.email = ["m@ku1ik.com"]
  s.homepage = "https://github.com/sickill/rainbow"
  s.licenses = ["MIT"]
  s.require_paths = ["lib"]
  s.required_ruby_version = Gem::Requirement.new(">= 1.9.2")
  s.rubygems_version = "2.0.14.1"
  s.summary = "Colorize printed text on ANSI terminals"

  if s.respond_to? :specification_version then
    s.specification_version = 4

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
      s.add_development_dependency(%q<bundler>, ["~> 1.3"])
    else
      s.add_dependency(%q<bundler>, ["~> 1.3"])
    end
  else
    s.add_dependency(%q<bundler>, ["~> 1.3"])
  end
end
