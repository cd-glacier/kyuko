ActiveRecord::Base.establish_connection(
  adapter:   'sqlite3',
  database:  'kyuko_info'
)

class Kyuko < ActiveRecord::Base
end


