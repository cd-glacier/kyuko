require 'active_record'
require './model.rb'

namespace :db do
  MIGRATIONS_DIR = 'db/migrate'

  # connect the database
  ActiveRecord::Base.establish_connection(
    :adapter => 'sqlite3',
    :database => 'database/db.sqlite'
  )

  desc "Migrate the database"
  task :migrate do
    ActiveRecord::Migrator.migrate(MIGRATIONS_DIR, ENV["VERSION"] ? ENV["VERSION"].to_i : nil)
  end

end


