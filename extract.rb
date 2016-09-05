
class Extraction

  def extract_time(line)
    time = line.split("I, [")[1].split("T")[0] 
  end

  def extract_kyuko(line)
    @kyuko = {}
    @kyuko[:when] = line.split("限目")[0]
    @kyuko[:name] = line.split(":")[1].split(" ")[0]
    @kyuko[:instructor] = line.split("講師(")[1].split(")")[0]
    return @kyuko
  end

  def get_tanabe()
    return @tanabe
  end

  def get_imade()
    return @imade
  end

  def extract(file)
    #0 = その日の情報とれていない,  1 = 田辺を取りたい, 
    #2 = 今出川を取りたい, 3 = その日の情報もうとった
    status = 0
    @imade= []
    @tanabe = []
    time = ""
    
    file.each_line do |line|
      if line.include?("I, [") then
        time = extract_time(line)
        status = 0
        next
      end

      if line.include?("田辺") && status == 0 then
        status = 1
        next
      elsif line.include?("今出川") && status == 1 then
        #田辺は取り終えた
        status = 2 
        next
      elsif line.include?("end") && status == 2 then
        stauts = 3
        next
      end

      if line.include?("限目") && status == 1 then
        tmp = extract_kyuko(line)
        tmp[:date] = time
        @tanabe << tmp 
      elsif line.include?("限目") && status == 2 then 
        tmp = extract_kyuko(line)
        tmp[:date] = time
        @imade << tmp 
      end
    end 
  end


end



File.open('./tmp/clockworkd.tweet.output') do |file|
  kyuko = Extraction.new
  kyuko.extract(file)

  #puts kyuko.get_tanabe()
end


