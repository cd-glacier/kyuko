# encoding: utf-8
# frozen_string_literal: true

module RuboCop
  module Cop
    module Lint
      # This cop checks for redundant access modifiers, including those with no
      # code, those which are repeated, and leading `public` modifiers in a
      # class or module body. Conditionally-defined methods are considered as
      # always being defined, and thus access modifiers guarding such methods
      # are not redundant.
      #
      # @example
      #
      #   class Foo
      #     public # this is redundant (default access is public)
      #
      #     def method
      #     end
      #
      #     private # this is not redundant (a method is defined)
      #     def method2
      #     end
      #
      #     private # this is redundant (no following methods are defined)
      #   end
      #
      # @example
      #
      #   class Foo
      #     # The following is not redundant (conditionally defined methods are
      #     # considered as always defining a method)
      #     private
      #
      #     if condition?
      #       def method
      #       end
      #     end
      #
      #     protected # this is not redundant (method is defined)
      #
      #     define_method(:method2) do
      #     end
      #
      #     protected # this is redundant (repeated from previous modifier)
      #
      #     [1,2,3].each do |i|
      #       define_method("foo#{i}") do
      #       end
      #     end
      #
      #     # The following is redundant (methods defined on the class'
      #     # singleton class are not affected by the public modifier)
      #     public
      #
      #     def self.method3
      #     end
      #   end
      class UselessAccessModifier < Cop
        MSG = 'Useless `%s` access modifier.'.freeze

        def on_class(node)
          check_node(node.children[2]) # class body
        end

        def on_module(node)
          check_node(node.children[1]) # module body
        end

        def on_block(node)
          return unless class_or_instance_eval?(node)

          check_node(node.children[2]) # block body
        end

        def on_sclass(node)
          check_node(node.children[1]) # singleton class body
        end

        private

        def_node_matcher :access_modifier, <<-PATTERN
          (send nil ${:public :protected :private})
        PATTERN

        def_node_matcher :static_method_definition?, <<-PATTERN
          {def (send nil {:attr :attr_reader :attr_writer :attr_accessor} ...)}
        PATTERN

        def_node_matcher :dynamic_method_definition?, <<-PATTERN
          {(send nil :define_method ...) (block (send nil :define_method ...) ...)}
        PATTERN

        def_node_matcher :class_or_instance_eval?, <<-PATTERN
          (block (send _ {:class_eval :instance_eval}) ...)
        PATTERN

        def check_node(node)
          return if node.nil?

          if node.begin_type?
            check_scope(node)
          elsif (vis = access_modifier(node))
            add_offense(node, :expression, format(MSG, vis))
          end
        end

        def check_scope(node)
          cur_vis, unused = check_child_nodes(node, nil, :public)

          add_offense(unused, :expression, format(MSG, cur_vis)) if unused
        end

        def check_child_nodes(node, unused, cur_vis)
          node.child_nodes.each do |child|
            if (new_vis = access_modifier(child))
              # does this modifier just repeat the existing visibility?
              if new_vis == cur_vis
                add_offense(child, :expression, format(MSG, cur_vis))
              else
                # was the previous modifier never applied to any defs?
                add_offense(unused, :expression, format(MSG, cur_vis)) if unused
                # once we have already warned about a certain modifier, don't
                # warn again even if it is never applied to any method defs
                unused = child
              end
              cur_vis = new_vis
            elsif method_definition?(child)
              unused = nil
            elsif start_of_new_scope?(child)
              check_scope(child)
            elsif !child.defs_type?
              cur_vis, unused = check_child_nodes(child, unused, cur_vis)
            end
          end

          [cur_vis, unused]
        end

        def method_definition?(child)
          static_method_definition?(child) || dynamic_method_definition?(child)
        end

        def start_of_new_scope?(child)
          child.module_type? || child.class_type? ||
            child.sclass_type? || class_or_instance_eval?(child)
        end
      end
    end
  end
end
