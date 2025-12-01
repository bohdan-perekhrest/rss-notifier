# frozen_string_literal: true

require 'async'
require 'async/http/internet'
require 'dotenv/load'

TELEGRAM_TOKEN = ENV.fetch('TELEGRAM_BOT_TOKEN')
CHAT_ID = ENV.fetch('TELEGRAM_CHAT_ID')

# Telegram sender adapter
class TelegramSender
  class << self
    def send_message(text)
      Async do
        internet = Async::HTTP::Internet.new
        internet.post(
          "https://api.telegram.org/bot#{TELEGRAM_TOKEN}/sendMessage",
          headers: { 'content-type' => 'application/json' },
          body: {
            chat_id: CHAT_ID,
            text: text,
            parse_mode: 'HTML',
            disable_web_page_preview: false,
            disable_notification: true
          }.to_json
        )
      ensure
        internet&.close
      end.wait
    end
  end
end
