# frozen_string_literal: true

require 'dotenv/load'
require 'httpx'

TELEGRAM_TOKEN = ENV.fetch('TELEGRAM_BOT_TOKEN')
CHAT_ID = ENV.fetch('TELEGRAM_CHAT_ID')

# Telegram sender adapter
class TelegramSender
  class << self
    def send_message(text)
      HTTPX
        .with(headers: { 'Content-Type': 'application/json' })
        .post(
          "https://api.telegram.org/bot#{TELEGRAM_TOKEN}/sendMessage",
          body: {
            chat_id: CHAT_ID,
            text: text,
            parse_mode: 'HTML',
            disable_web_page_preview: false,
            disable_notification: true
          }.to_json
        )
    end
  end
end
