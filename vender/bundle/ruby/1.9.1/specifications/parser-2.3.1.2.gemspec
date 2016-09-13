# -*- encoding: utf-8 -*-

Gem::Specification.new do |s|
  s.name = "parser"
  s.version = "2.3.1.2"

  s.required_rubygems_version = Gem::Requirement.new(">= 0") if s.respond_to? :required_rubygems_version=
  s.authors = ["whitequark"]
  s.date = "2016-06-02"
  s.description = "A Ruby parser written in pure Ruby."
  s.email = ["whitequark@whitequark.org"]
  s.executables = ["ruby-parse", "ruby-rewrite"]
  s.files = ["bin/ruby-parse", "bin/ruby-rewrite"]
  s.homepage = "https://github.com/whitequark/parser"
  s.licenses = ["MIT"]
  s.require_paths = ["lib"]
  s.rubygems_version = "1.8.23"
  s.summary = "A Ruby parser written in pure Ruby."

  if s.respond_to? :specification_version then
    s.specification_version = 4

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
      s.add_runtime_dependency(%q<ast>, ["~> 2.2"])
      s.add_development_dependency(%q<bundler>, ["~> 1.2"])
      s.add_development_dependency(%q<rake>, ["~> 10.0"])
      s.add_development_dependency(%q<racc>, ["= 1.4.14"])
      s.add_development_dependency(%q<cliver>, ["~> 0.3.0"])
      s.add_development_dependency(%q<yard>, [">= 0"])
      s.add_development_dependency(%q<kramdown>, [">= 0"])
      s.add_development_dependency(%q<minitest>, ["~> 5.0"])
      s.add_development_dependency(%q<simplecov>, ["~> 0.8.2"])
      s.add_development_dependency(%q<coveralls>, [">= 0"])
      s.add_development_dependency(%q<json_pure>, [">= 0"])
      s.add_development_dependency(%q<mime-types>, ["~> 1.25"])
      s.add_development_dependency(%q<rest-client>, ["~> 1.6.7"])
      s.add_development_dependency(%q<simplecov-sublime-ruby-coverage>, [">= 0"])
      s.add_development_dependency(%q<gauntlet>, [">= 0"])
    else
      s.add_dependency(%q<ast>, ["~> 2.2"])
      s.add_dependency(%q<bundler>, ["~> 1.2"])
      s.add_dependency(%q<rake>, ["~> 10.0"])
      s.add_dependency(%q<racc>, ["= 1.4.14"])
      s.add_dependency(%q<cliver>, ["~> 0.3.0"])
      s.add_dependency(%q<yard>, [">= 0"])
      s.add_dependency(%q<kramdown>, [">= 0"])
      s.add_dependency(%q<minitest>, ["~> 5.0"])
      s.add_dependency(%q<simplecov>, ["~> 0.8.2"])
      s.add_dependency(%q<coveralls>, [">= 0"])
      s.add_dependency(%q<json_pure>, [">= 0"])
      s.add_dependency(%q<mime-types>, ["~> 1.25"])
      s.add_dependency(%q<rest-client>, ["~> 1.6.7"])
      s.add_dependency(%q<simplecov-sublime-ruby-coverage>, [">= 0"])
      s.add_dependency(%q<gauntlet>, [">= 0"])
    end
  else
    s.add_dependency(%q<ast>, ["~> 2.2"])
    s.add_dependency(%q<bundler>, ["~> 1.2"])
    s.add_dependency(%q<rake>, ["~> 10.0"])
    s.add_dependency(%q<racc>, ["= 1.4.14"])
    s.add_dependency(%q<cliver>, ["~> 0.3.0"])
    s.add_dependency(%q<yard>, [">= 0"])
    s.add_dependency(%q<kramdown>, [">= 0"])
    s.add_dependency(%q<minitest>, ["~> 5.0"])
    s.add_dependency(%q<simplecov>, ["~> 0.8.2"])
    s.add_dependency(%q<coveralls>, [">= 0"])
    s.add_dependency(%q<json_pure>, [">= 0"])
    s.add_dependency(%q<mime-types>, ["~> 1.25"])
    s.add_dependency(%q<rest-client>, ["~> 1.6.7"])
    s.add_dependency(%q<simplecov-sublime-ruby-coverage>, [">= 0"])
    s.add_dependency(%q<gauntlet>, [">= 0"])
  end
end
