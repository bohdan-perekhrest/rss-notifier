# frozen_string_literal: true

require_relative 'base'

# Atom RSS reader - handles standard Atom feeds
class AtomReader < BaseReader
  private

  def extract_channel_name(feed)
    feed.channel.title
  end

  def each_entry(feed, &block)
    feed.channel.items.each(&block)
  end

  def extract_url(entry)
    entry.link
  end

  def extract_title(entry)
    entry.title
  end

  def should_skip?(entry, _url, last_check_at)
    return true if entry.pubDate < (Time.now - 86400) && last_check_at.nil?

    last_check_at && entry.pubDate < last_check_at
  end
end
