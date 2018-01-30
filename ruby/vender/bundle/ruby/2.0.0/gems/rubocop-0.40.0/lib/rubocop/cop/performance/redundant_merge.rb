# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Performance
      # This cop identifies places where `Hash#merge!` can be replaced by
      # `Hash#[]=`.
      #
      # @example
      #   hash.merge!(a: 1)
      #   hash.merge!({'key' => 'value'})
      #   hash.merge!(a: 1, b: 2)
      class RedundantMerge < Cop
        AREF_ASGN = '%s[%s] = %s'.freeze
        MSG = 'Use `%s` instead of `%s`.'.freeze

        def_node_matcher :redundant_merge, '(send $_ :merge! (hash $...))'
        def_node_matcher :modifier_flow_control, '[{if while until} #modifier?]'
        def_node_matcher :each_with_object_node, <<-END
          (block (send _ :each_with_object _) (args _ $_) ...)
        END

        def on_send(node)
          each_redundant_merge(node) do |receiver, pairs|
            assignments = to_assignments(receiver, pairs).join('; ')
            message = format(MSG, assignments, node.source)
            add_offense(node, :expression, message)
          end
        end

        def autocorrect(node)
          redundant_merge(node) do |receiver, pairs|
            lambda do |corrector|
              new_source = to_assignments(receiver, pairs).join("\n")

              parent = node.parent
              if parent && pairs.size > 1
                if modifier_flow_control(parent)
                  new_source = rewrite_with_modifier(node, parent, new_source)
                  node = parent
                else
                  padding = "\n#{leading_spaces(node)}"
                  new_source.gsub!(/\n/, padding)
                end
              end

              corrector.replace(node.source_range, new_source)
            end
          end
        end

        private

        def each_redundant_merge(node)
          redundant_merge(node) do |receiver, pairs|
            next if node.value_used? &&
                    !value_used_inside_each_with_object?(node, receiver)
            next if pairs.size > 1 && !receiver.pure?
            next if pairs.size > max_key_value_pairs

            yield receiver, pairs
          end
        end

        def value_used_inside_each_with_object?(node, receiver)
          while receiver.respond_to?(:send_type?) && receiver.send_type?
            receiver, = *receiver
          end

          unless receiver.respond_to?(:lvar_type?) && receiver.lvar_type?
            return false
          end

          parent = node.parent
          grandparent = parent.parent if parent.begin_type?
          second_arg = each_with_object_node(grandparent || parent)
          return false if second_arg.nil?

          receiver.loc.name.source == second_arg.loc.name.source
        end

        def to_assignments(receiver, pairs)
          pairs.map do |pair|
            key, value = *pair
            key_src = if key.sym_type? && !key.source.start_with?(':')
                        ":#{key.source}"
                      else
                        key.source
                      end

            format(AREF_ASGN, receiver.source, key_src, value.source)
          end
        end

        def rewrite_with_modifier(node, parent, new_source)
          cond, = *parent
          padding = "\n#{(' ' * indent_width) + leading_spaces(node)}"
          new_source.gsub!(/\n/, padding)

          parent.loc.keyword.source << ' ' << cond.source << padding <<
            new_source << "\n" << leading_spaces(node) << 'end'
        end

        def leading_spaces(node)
          node.source_range.source_line[/\A\s*/]
        end

        def indent_width
          @config.for_cop('IndentationWidth')['Width'] || 2
        end

        def modifier?(node)
          node.loc.respond_to?(:end) && node.loc.end.nil?
        end

        def max_key_value_pairs
          cop_config['MaxKeyValuePairs'].to_i
        end
      end
    end
  end
end
