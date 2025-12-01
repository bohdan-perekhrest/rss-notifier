# frozen_string_literal: true

# Cache adapter
class Cache
  CACHE_FILE = 'cache'

  class << self
    def read_last_check
      return nil unless File.exist?(CACHE_FILE)

      Time.parse(File.read(CACHE_FILE).strip)
    rescue ArgumentError => e
      warn "Failed to parse cache: #{e.message}"
      nil
    end

    def write_last_check(time = Time.now)
      File.write(CACHE_FILE, time.to_s)
    end
  end
end
