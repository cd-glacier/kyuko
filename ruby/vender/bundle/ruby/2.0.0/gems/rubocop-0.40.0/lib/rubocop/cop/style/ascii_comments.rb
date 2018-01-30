# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Style
      # This cop checks for non-ascii (non-English) characters
      # in comments.
      class AsciiComments < Cop
        MSG = 'Use only ascii symbols in comments.'.freeze

        def investigate(processed_source)
          processed_source.comments.each do |comment|
            add_offense(comment, :expression) unless comment.text.ascii_only?
          end
        end
      end
    end
  end
end
