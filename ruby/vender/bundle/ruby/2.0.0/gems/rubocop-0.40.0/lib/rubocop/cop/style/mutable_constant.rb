# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Style
      # This cop checks whether some constant value isn't a
      # mutable literal (e.g. array or hash).
      #
      # @example
      #   # bad
      #   CONST = [1, 2, 3]
      #
      #   # good
      #   CONST = [1, 2, 3].freeze
      class MutableConstant < Cop
        MSG = 'Freeze mutable objects assigned to constants.'.freeze

        include FrozenStringLiteral

        def on_casgn(node)
          _scope, _const_name, value = *node
          on_assignment(value)
        end

        def on_or_asgn(node)
          lhs, value = *node
          on_assignment(value) if lhs && lhs.type == :casgn
        end

        def autocorrect(node)
          expr = node.source_range
          ->(corrector) { corrector.replace(expr, "#{expr.source}.freeze") }
        end

        private

        def on_assignment(value)
          return unless value
          return unless value.mutable_literal?
          return if FROZEN_STRING_LITERAL_TYPES.include?(value.type) &&
                    frozen_string_literals_enabled?(processed_source)

          add_offense(value, :expression)
        end
      end
    end
  end
end
