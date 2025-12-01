# frozen_string_literal: true

require 'sqlite3'

# DB adapter class
class DB
  @instance = nil
  @mutex = Mutex.new

  class << self
    def instance
      @mutex.synchronize do
        @instance ||= new
      end
    end

    def video_seen?(url)
      instance.video_seen?(url)
    end

    def mark_video_seen(url, title)
      instance.mark_video_seen(url, title)
    end

    def close
      @instance&.db&.close
      @instance = nil
    end
  end

  attr_reader :db

  def initialize
    @db = SQLite3::Database.new('videos.db')
    @db.execute <<-SQL
      CREATE TABLE IF NOT EXISTS videos (
        url TEXT PRIMARY KEY,
        title TEXT,
        seen_at DATETIME DEFAULT CURRENT_TIMESTAMP
      );
    SQL
  end

  def video_seen?(url)
    result = db.execute('SELECT 1 FROM videos WHERE url = ?', url)
    !result.empty?
  end

  def mark_video_seen(url, title)
    db.execute('INSERT OR IGNORE INTO videos (url, title) VALUES (?, ?)', [url, title])
  end
end
