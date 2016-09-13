# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Style
      # This cops checks if empty lines around the bodies of modules match
      # the configuration.
      #
      # @example
      #
      #   module Test
      #
      #      def something
      #        ...
      #      end
      #
      #   end
      #
      class EmptyLinesAroundModuleBody < Cop
        include EmptyLinesAroundBody

        KIND = 'module'.freeze

        def on_module(node)
          _name, body = *node
          check(node, body)
        end
      end
    end
  end
end
