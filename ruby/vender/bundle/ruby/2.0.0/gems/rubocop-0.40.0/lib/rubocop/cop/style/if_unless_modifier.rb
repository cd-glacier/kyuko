# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Style
      # Checks for if and unless statements that would fit on one line
      # if written as a modifier if/unless.
      # The maximum line length is configurable.
      class IfUnlessModifier < Cop
        include IfNode
        include StatementModifier

        ASSIGNMENT_TYPES = [:lvasgn, :casgn, :cvasgn,
                            :gvasgn, :ivasgn, :masgn].freeze

        def message(keyword)
          "Favor modifier `#{keyword}` usage when having a single-line body." \
          ' Another good alternative is the usage of control flow `&&`/`||`.'
        end

        def on_if(node)
          # discard ternary ops, if/else and modifier if/unless nodes
          return if ternary?(node)
          return if modifier_if?(node)
          return if elsif?(node)
          return if if_else?(node)
          return if node.chained?
          return unless fit_within_line_as_modifier_form?(node)
          return if nested_conditional?(node)
          add_offense(node, :keyword, message(node.loc.keyword.source))
        end

        def parenthesize?(node)
          # Parenthesize corrected expression if changing to modifier-if form
          # would change the meaning of the parent expression
          # (due to the low operator precedence of modifier-if)
          return false if node.parent.nil?
          return true if ASSIGNMENT_TYPES.include?(node.parent.type)

          if node.parent.send_type?
            _receiver, _name, *args = *node.parent
            return !method_uses_parens?(node.parent, args.first)
          end

          false
        end

        def method_uses_parens?(node, limit)
          source = node.source_range.source_line[0...limit.loc.column]
          source =~ /\s*\(\s*$/
        end

        def autocorrect(node)
          cond, body, _else = if_node_parts(node)

          oneline =
            "#{body.source} #{node.loc.keyword.source} " + cond.source
          first_line_comment = processed_source.comments.find do |c|
            c.loc.line == node.loc.line
          end
          if first_line_comment
            oneline << ' ' << first_line_comment.loc.expression.source
          end
          oneline = "(#{oneline})" if parenthesize?(node)

          ->(corrector) { corrector.replace(node.source_range, oneline) }
        end

        private

        # returns false if the then or else children are conditionals
        def nested_conditional?(node)
          node.children[1, 2].any? { |child| child && child.type == :if }
        end
      end
    end
  end
end
