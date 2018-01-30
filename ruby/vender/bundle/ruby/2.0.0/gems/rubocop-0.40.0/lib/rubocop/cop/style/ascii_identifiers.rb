# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Style
      # This cop checks for non-ascii characters in identifier names.
      class AsciiIdentifiers < Cop
        MSG = 'Use only ascii symbols in identifiers.'.freeze

        def investigate(processed_source)
          processed_source.tokens.each do |t|
            next unless t.type == :tIDENTIFIER && !t.text.ascii_only?
            add_offense(nil, t.pos)
          end
        end
      end
    end
  end
end
