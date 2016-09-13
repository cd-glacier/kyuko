# -*- encoding: utf-8 -*-

Gem::Specification.new do |s|
  s.name = "daemons"
  s.version = "1.2.3"

  s.required_rubygems_version = Gem::Requirement.new(">= 0") if s.respond_to? :required_rubygems_version=
  s.authors = ["Thomas Uehlinger"]
  s.date = "2015-03-08"
  s.description = "    Daemons provides an easy way to wrap existing ruby scripts (for example a\n    self-written server)  to be run as a daemon and to be controlled by simple\n    start/stop/restart commands.\n\n    You can also call blocks as daemons and control them from the parent or just\n    daemonize the current process.\n\n    Besides this basic functionality, daemons offers many advanced features like\n    exception backtracing and logging (in case your ruby script crashes) and\n    monitoring and automatic restarting of your processes if they crash.\n"
  s.email = "thomas.uehinger@gmail.com"
  s.homepage = "https://github.com/thuehlinger/daemons"
  s.licenses = ["MIT"]
  s.require_paths = ["lib"]
  s.rubygems_version = "1.8.23"
  s.summary = "A toolkit to create and control daemons in different ways"

  if s.respond_to? :specification_version then
    s.specification_version = 4

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
    else
    end
  else
  end
end
