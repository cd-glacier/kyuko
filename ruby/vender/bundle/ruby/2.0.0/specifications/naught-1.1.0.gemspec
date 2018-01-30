# -*- encoding: utf-8 -*-

Gem::Specification.new do |s|
  s.name = "naught"
  s.version = "1.1.0"

  s.required_rubygems_version = Gem::Requirement.new(">= 0") if s.respond_to? :required_rubygems_version=
  s.authors = ["Avdi Grimm"]
  s.date = "2015-09-08"
  s.description = "Naught is a toolkit for building Null Objects"
  s.email = ["avdi@avdi.org"]
  s.homepage = "https://github.com/avdi/naught"
  s.licenses = ["MIT"]
  s.require_paths = ["lib"]
  s.rubygems_version = "2.0.14.1"
  s.summary = "Naught is a toolkit for building Null Objects"

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
