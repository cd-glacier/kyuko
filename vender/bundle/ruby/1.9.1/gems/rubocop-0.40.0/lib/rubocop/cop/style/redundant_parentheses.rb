# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Style
      # This cop checks for redundant parentheses.
      #
      # @example
      #
      #   # bad
      #   (x) if ((y.z).nil?)
      #
      #   # good
      #   x if y.z.nil?
      #
      class RedundantParentheses < Cop
        include Parentheses

        ALLOWED_LITERALS = [:irange, :erange].freeze

        def_node_matcher :square_brackets?, '(send (send _recv _msg) :[] ...)'
        def_node_matcher :range_end?, '^^{irange erange}'
        def_node_matcher :method_node_and_args, '$(send _recv _msg $...)'
        def_node_matcher :rescue?, '{^resbody ^^resbody}'
        def_node_matcher :arg_in_call_with_block?,
                         '^^(block (send _ _ equal?(%0) ...) ...)'

        def on_begin(node)
          return if !parentheses?(node) || parens_allowed?(node)
          check(node)
        end

        def parens_allowed?(node)
          child  = node.children.first
          parent = node.parent

          # don't flag `break(1)`, etc
          (keyword_ancestor?(node) && parens_required?(node)) ||
            # don't flag `method ({key: value})`
            (child.hash_type? && first_arg?(node) && !parentheses?(parent)) ||
            # don't flag `rescue(ExceptionClass)`
            rescue?(node) ||
            # don't flag `method (arg) { }`
            (arg_in_call_with_block?(node) && !parentheses?(parent)) ||
            # don't flag
            # ```
            # { a: (1
            #      ), }
            # ```
            allowed_array_or_hash_element?(node)
        end

        def check(begin_node)
          node = begin_node.children.first
          if keyword_with_redundant_parentheses?(node)
            return offense(begin_node, 'a keyword')
          end
          return offense(begin_node, 'a literal') if disallowed_literal?(node)
          return offense(begin_node, 'a variable') if node.variable?
          return offense(begin_node, 'a constant') if node.const_type?
          check_send(begin_node, node) if node.send_type?
        end

        def check_send(begin_node, node)
          if node.unary_operation?
            return if begin_node.chained?

            # parens are not redundant in `(!recv.method arg)`
            node = node.children.first while node.unary_operation?
            if node.send_type?
              return unless method_call_with_redundant_parentheses?(node)
            end

            offense(begin_node, 'an unary operation')
          else
            return unless method_call_with_redundant_parentheses?(node)
            return if call_chain_starts_with_int?(begin_node, node)

            offense(begin_node, 'a method call')
          end
        end

        def offense(node, msg)
          add_offense(node, :expression, "Don't use parentheses around #{msg}.")
        end

        def keyword_ancestor?(node)
          node.parent && node.parent.keyword?
        end

        def allowed_array_or_hash_element?(node)
          (hash_element?(node) || array_element?(node)) &&
            only_closing_paren_before_comma?(node)
        end

        def hash_element?(node)
          node.parent && node.parent.pair_type?
        end

        def array_element?(node)
          node.parent && node.parent.array_type?
        end

        def only_closing_paren_before_comma?(node)
          source_buffer = node.source_range.source_buffer
          line_range = source_buffer.line_range(node.loc.end.line)

          line_range.source =~ /^\s*\)\s*,/
        end

        def disallowed_literal?(node)
          node.literal? && !ALLOWED_LITERALS.include?(node.type)
        end

        def keyword_with_redundant_parentheses?(node)
          return false unless node.keyword?
          return true if node.special_keyword?

          args = *node

          if args.size == 1 && args.first && args.first.begin_type?
            parentheses?(args.first)
          else
            args.empty? || parentheses?(node)
          end
        end

        def method_call_with_redundant_parentheses?(node)
          return false unless node.send_type?
          return false if node.keyword_not?
          return false if range_end?(node)

          send_node, args = method_node_and_args(node)

          args.empty? || parentheses?(send_node) || square_brackets?(send_node)
        end

        def first_arg?(node)
          send_node = node.parent
          return false unless send_node && send_node.send_type?

          _receiver, _method_name, *args = *send_node
          node.equal?(args.first)
        end

        def call_chain_starts_with_int?(begin_node, send_node)
          recv = first_part_of_call_chain(send_node)
          recv && recv.int_type? && (parent = begin_node.parent) &&
            parent.send_type? &&
            (parent.method_name == :-@ || parent.method_name == :+@)
        end
      end
    end
  end
end
