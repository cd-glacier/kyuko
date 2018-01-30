# encoding: utf-8
# frozen_string_literal: true

require 'set'

module RuboCop
  module Cop
    module Style
      # This cop checks for extra/unnecessary whitespace.
      #
      # @example
      #
      #   # good if AllowForAlignment is true
      #   name      = "RuboCop"
      #   # Some comment and an empty line
      #
      #   website  += "/bbatsov/rubocop" unless cond
      #   puts        "rubocop"          if     debug
      #
      #   # bad for any configuration
      #   set_app("RuboCop")
      #   website  = "https://github.com/bbatsov/rubocop"
      class ExtraSpacing < Cop
        include PrecedingFollowingAlignment

        MSG_UNNECESSARY = 'Unnecessary spacing detected.'.freeze
        MSG_UNALIGNED_ASGN = '`=` is not aligned with the %s assignment.'.freeze

        def investigate(processed_source)
          if force_equal_sign_alignment?
            @asgn_tokens = processed_source.tokens.select { |t| equal_sign?(t) }
            # Only attempt to align the first = on each line
            @asgn_tokens = Set.new(@asgn_tokens.uniq { |t| t.pos.line })
            @asgn_lines  = @asgn_tokens.map { |t| t.pos.line }
            # Don't attempt to correct the same = more than once
            @corrected   = Set.new
          end

          processed_source.tokens.each_cons(2) do |t1, t2|
            check_tokens(processed_source.ast, t1, t2)
          end
        end

        def autocorrect(range)
          lambda do |corrector|
            if range.source.end_with?('=')
              align_equal_sign(range, corrector)
            else
              corrector.remove(range)
            end
          end
        end

        private

        def check_tokens(ast, t1, t2)
          return if t2.type == :tNL

          if force_equal_sign_alignment? &&
             @asgn_tokens.include?(t2) &&
             (@asgn_lines.include?(t2.pos.line - 1) ||
              @asgn_lines.include?(t2.pos.line + 1))
            check_assignment(t2)
          else
            check_other(t1, t2, ast)
          end
        end

        def check_assignment(token)
          assignment_line = ''
          message = ''
          if should_aligned_with_preceding_line?(token)
            assignment_line = preceding_line(token)
            message = format(MSG_UNALIGNED_ASGN, 'preceding')
          else
            assignment_line = following_line(token)
            message = format(MSG_UNALIGNED_ASGN, 'following')
          end
          return if aligned_assignment?(token.pos, assignment_line)
          add_offense(token.pos, token.pos, message)
        end

        def should_aligned_with_preceding_line?(token)
          @asgn_lines.include?(token.pos.line - 1)
        end

        def preceding_line(token)
          processed_source.lines[token.pos.line - 2]
        end

        def following_line(token)
          processed_source.lines[token.pos.line]
        end

        def check_other(t1, t2, ast)
          return if t1.pos.line != t2.pos.line
          return if t2.pos.begin_pos - 1 <= t1.pos.end_pos
          return if allow_for_alignment? && aligned_tok?(t2)

          start_pos = t1.pos.end_pos
          return if ignored_ranges(ast).find { |r| r.include?(start_pos) }

          end_pos = t2.pos.begin_pos - 1
          range = Parser::Source::Range.new(processed_source.buffer,
                                            start_pos, end_pos)
          # Unary + doesn't appear as a token and needs special handling.
          return if unary_plus_non_offense?(range)

          add_offense(range, range, MSG_UNNECESSARY)
        end

        def aligned_tok?(token)
          if token.type == :tCOMMENT
            aligned_comments?(token)
          else
            aligned_with_something?(token.pos)
          end
        end

        def unary_plus_non_offense?(range)
          range.resize(range.size + 1).source =~ /^ ?\+$/
        end

        # Returns an array of ranges that should not be reported. It's the
        # extra spaces between the keys and values in a hash, since those are
        # handled by the Style/AlignHash cop.
        def ignored_ranges(ast)
          return [] unless ast

          @ignored_ranges ||= on_node(:pair, ast).map do |pair|
            key, value = *pair
            key.source_range.end_pos...value.source_range.begin_pos
          end
        end

        def aligned_comments?(token)
          ix = processed_source.comments.index do |c|
            c.loc.expression.begin_pos == token.pos.begin_pos
          end
          aligned_with_previous_comment?(ix) || aligned_with_next_comment?(ix)
        end

        def aligned_with_previous_comment?(ix)
          ix > 0 && comment_column(ix - 1) == comment_column(ix)
        end

        def aligned_with_next_comment?(ix)
          ix < processed_source.comments.length - 1 &&
            comment_column(ix + 1) == comment_column(ix)
        end

        def comment_column(ix)
          processed_source.comments[ix].loc.column
        end

        def force_equal_sign_alignment?
          cop_config['ForceEqualSignAlignment']
        end

        def equal_sign?(token)
          token.type == :tEQL || token.type == :tOP_ASGN
        end

        def align_equal_sign(range, corrector)
          lines  = contiguous_assignment_lines(range)
          tokens = @asgn_tokens.select { |t| lines.include?(t.pos.line) }

          columns  = tokens.map { |t| align_column(t) }
          align_to = columns.max

          tokens.each do |token|
            next unless @corrected.add?(token)
            diff = align_to - token.pos.last_column

            if diff > 0
              corrector.insert_before(token.pos, ' ' * diff)
            elsif diff < 0
              corrector.remove_preceding(token.pos, -diff)
            end
          end
        end

        def contiguous_assignment_lines(range)
          result = [range.line]

          range.line.downto(1) do |lineno|
            @asgn_lines.include?(lineno) ? result << lineno : break
          end
          range.line.upto(processed_source.lines.size) do |lineno|
            @asgn_lines.include?(lineno) ? result << lineno : break
          end

          result.sort!
        end

        def align_column(asgn_token)
          # if we removed unneeded spaces from the beginning of this =,
          # what column would it end from?
          line    = processed_source.lines[asgn_token.pos.line - 1]
          leading = line[0...asgn_token.pos.column]
          spaces  = leading.size - (leading =~ / *\Z/)
          asgn_token.pos.last_column - spaces + 1
        end
      end
    end
  end
end
