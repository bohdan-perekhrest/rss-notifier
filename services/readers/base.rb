# frozen_string_literal: true

require 'async'
require 'async/http/internet'
require 'rss'
require_relative '../telegram_sender'

# Base RSS Reader with Template Method pattern
class BaseReader
  def self.parse(url, last_check_at)
    new.parse(url, last_check_at)
  end

  def parse(url, last_check_at)
    feed = get_feed(url)
    channel_name = extract_channel_name(feed)

    each_entry(feed) do |entry|
      url = extract_url(entry)
      next if should_skip?(entry, url, last_check_at)

      title = extract_title(entry)
      message = format_message(channel_name, url, title)

      TelegramSender.send_message(message)
    end
  end

  private

  def get_feed(url)
    Async do
      internet = Async::HTTP::Internet.new
      response = internet.get(url)
      RSS::Parser.parse(response.read, validate: false)
    ensure
      internet&.close
    end.wait
  end

  # Template methods to be implemented by subclasses
  def extract_channel_name(feed)
    raise NotImplementedError, 'subclasses must implement extract_channel_name'
  end

  def each_entry(feed, &block)
    raise NotImplementedError, 'subclasses must implement each_entry'
  end

  def extract_url(entry)
    raise NotImplementedError, 'subclasses must implement extract_url'
  end

  def extract_title(entry)
    raise NotImplementedError, 'subclasses must implement extract_title'
  end

  # Hook method - can be overridden by subclasses
  def should_skip?(_entry, _url, _last_check_at)
    false
  end

  def format_message(channel_name, url, title)
    <<~MSG
      FROM: <b>#{escape_html(channel_name)}</b>
      <a href="#{url}">#{escape_html(title)}</a>
    MSG
  end

  def escape_html(text)
    text.to_s.gsub('&', '&amp;').gsub('<', '&lt;').gsub('>', '&gt;')
  end
end
