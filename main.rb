# frozen_string_literal: true

require_relative 'services/db'
require_relative 'services/telegram_sender'
require 'rss'
require 'yaml'

channels = YAML.load_file('config/channels.yml')['channels']
list = channels.map { |id| "https://www.youtube.com/feeds/videos.xml?channel_id=#{id}" }

def escape_html(text)
  text
    .to_s
    .gsub('&', '&amp;')
    .gsub('<', '&lt;')
    .gsub('>', '&gt;')
end

list.each do |feed_url|
  feed = RSS::Parser.parse(feed_url, validate: false)
  channel_name = feed.author.name.content

  feed.entries.each do |entry|
    video_url = entry.link.href
    next if video_url.include?('shorts')
    next if DB.video_seen?(video_url)

    video_title = entry.title.content

    message = <<~MSG
      Author: <b>#{escape_html(channel_name)}</b>
      Published at: #{entry.published.strftime('%Y-%m-%d %H:%M')}

      <a href="#{video_url}">#{escape_html(video_title)}</a>
    MSG

    TelegramSender.send_message(message)
    DB.mark_video_seen(video_url, video_title)
  end
end
