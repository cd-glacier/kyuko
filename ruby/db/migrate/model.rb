require "active_record"

ActiveRecord::Base.establish_connection(
  adapter:   'sqlite3',
  database:  'kyuko_info'
)

class InitialSchema < ActiveRecord::Migration
  def self.up
    create_table :kyukos do |t|
      t.string :name
      t.integer :when
      t.string :instructor
      t.datetime :date
      t.timestamp
    end
  end

  def self.down
    drop_table :users
  end
end

=begin
class Kyuko < ActiveRecord::Base
end
=end

