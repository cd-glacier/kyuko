# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Lint
      # This cop checks for *rescue* blocks targeting the Exception class.
      class RescueException < Cop
        MSG = 'Avoid rescuing the `Exception` class. ' \
              'Perhaps you meant to rescue `StandardError`?'.freeze

        def on_resbody(node)
          return unless node.children.first
          rescue_args = node.children.first.children
          return unless rescue_args.any? { |a| targets_exception?(a) }

          add_offense(node, :expression)
        end

        def targets_exception?(rescue_arg_node)
          rescue_arg_node.const_name == 'Exception'
        end
      end
    end
  end
end
