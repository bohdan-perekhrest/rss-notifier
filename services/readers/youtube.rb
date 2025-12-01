# frozen_string_literal: true

require_relative 'base'

# Youtube RSS reader - handles YouTube Atom feeds
class YoutubeReader < BaseReader
  private

  def extract_channel_name(feed)
    feed.author.name.content
  end

  def each_entry(feed, &block)
    feed.entries.each(&block)
  end

  def extract_url(entry)
    entry.link.href
  end

  def extract_title(entry)
    entry.title.content
  end

  def should_skip?(entry, url, last_check_at)
    return true if entry.published.content < (Time.now - 86400) && last_check_at.nil?

    url.include?('shorts') || (last_check_at && entry.published.content < last_check_at)
  end
end
