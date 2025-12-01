# frozen_string_literal: true

require_relative 'services/readers/youtube'
require_relative 'services/readers/atom'
require_relative 'services/cache'
require 'async'
require 'yaml'

config = YAML.load_file('config/feeds.yml')
feeds = config['feeds'] || []
last_check = Cache.read_last_check

Async do
  feeds.map do |feed|
    Async do
      case feed['type']
      when 'youtube'
        url = "https://www.youtube.com/feeds/videos.xml?channel_id=#{feed['id']}"
        YoutubeReader.parse(url, last_check)
      when 'rss'
        AtomReader.parse(feed['url'], last_check)
      else
        warn "Unknown feed type: #{feed['type']}"
      end
    end
  end.map(&:wait)
end

Cache.write_last_check
