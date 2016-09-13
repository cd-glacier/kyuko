# -*- encoding: utf-8 -*-

Gem::Specification.new do |s|
  s.name = "twitter"
  s.version = "5.16.0"

  s.required_rubygems_version = Gem::Requirement.new(">= 1.3.5") if s.respond_to? :required_rubygems_version=
  s.authors = ["Erik Michaels-Ober", "John Nunemaker", "Wynn Netherland", "Steve Richert", "Steve Agalloco"]
  s.date = "2016-01-24"
  s.description = "A Ruby interface to the Twitter API."
  s.email = ["sferik@gmail.com"]
  s.homepage = "http://sferik.github.com/twitter/"
  s.licenses = ["MIT"]
  s.require_paths = ["lib"]
  s.rubygems_version = "1.8.23"
  s.summary = "A Ruby interface to the Twitter API."

  if s.respond_to? :specification_version then
    s.specification_version = 4

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
      s.add_runtime_dependency(%q<addressable>, ["~> 2.3"])
      s.add_runtime_dependency(%q<buftok>, ["~> 0.2.0"])
      s.add_runtime_dependency(%q<equalizer>, ["= 0.0.10"])
      s.add_runtime_dependency(%q<faraday>, ["~> 0.9.0"])
      s.add_runtime_dependency(%q<http>, ["~> 1.0"])
      s.add_runtime_dependency(%q<http_parser.rb>, ["~> 0.6.0"])
      s.add_runtime_dependency(%q<json>, ["~> 1.8"])
      s.add_runtime_dependency(%q<memoizable>, ["~> 0.4.0"])
      s.add_runtime_dependency(%q<naught>, ["~> 1.0"])
      s.add_runtime_dependency(%q<simple_oauth>, ["~> 0.3.0"])
      s.add_development_dependency(%q<bundler>, ["~> 1.0"])
    else
      s.add_dependency(%q<addressable>, ["~> 2.3"])
      s.add_dependency(%q<buftok>, ["~> 0.2.0"])
      s.add_dependency(%q<equalizer>, ["= 0.0.10"])
      s.add_dependency(%q<faraday>, ["~> 0.9.0"])
      s.add_dependency(%q<http>, ["~> 1.0"])
      s.add_dependency(%q<http_parser.rb>, ["~> 0.6.0"])
      s.add_dependency(%q<json>, ["~> 1.8"])
      s.add_dependency(%q<memoizable>, ["~> 0.4.0"])
      s.add_dependency(%q<naught>, ["~> 1.0"])
      s.add_dependency(%q<simple_oauth>, ["~> 0.3.0"])
      s.add_dependency(%q<bundler>, ["~> 1.0"])
    end
  else
    s.add_dependency(%q<addressable>, ["~> 2.3"])
    s.add_dependency(%q<buftok>, ["~> 0.2.0"])
    s.add_dependency(%q<equalizer>, ["= 0.0.10"])
    s.add_dependency(%q<faraday>, ["~> 0.9.0"])
    s.add_dependency(%q<http>, ["~> 1.0"])
    s.add_dependency(%q<http_parser.rb>, ["~> 0.6.0"])
    s.add_dependency(%q<json>, ["~> 1.8"])
    s.add_dependency(%q<memoizable>, ["~> 0.4.0"])
    s.add_dependency(%q<naught>, ["~> 1.0"])
    s.add_dependency(%q<simple_oauth>, ["~> 0.3.0"])
    s.add_dependency(%q<bundler>, ["~> 1.0"])
  end
end
