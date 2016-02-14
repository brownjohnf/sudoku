#!/usr/bin/env ruby

module Sudoku
  class CLI
    class << self
      # Internal: Call this to use the CLI interface.
      #
      # Returns a new Board, instantiated from ARGF.
      def main
        board = ARGF.read.gsub(/\n/, ',').gsub("-", "0").split(',').map(&:to_i)

        Board.new(board)
      end
    end
  end

  class Board
    # All the valid numbers to try in our board. 0 means unset.
    NUMS = (0..9).to_a.reverse

    # Public: Create a new board.
    #
    # board   - Array of 81 Fixnum values, which is the rows of the board
    #           joined end-to-end.
    #
    # Returns a new Board.
    def initialize(board)
      @board = board
      raise ArgumentError, "board must be 9x9!" unless @board.length == 81
    end

    # Public: Print the board in an attractive way to STDOUT.
    #
    # Returns nil.
    def display
      @board.each_slice(9) do |row|
        puts row.join(",")
      end
    end

    # Public: Check to see whether the board is in a valid state.
    #
    # Returns Boolean whether or not the board is in a valid state.
    def valid?
      81.times do |i|
        return false if @board[i] == 0
        return false unless check?(i)
      end

      true
    end

    # Public: Solve the board, from its current state.
    #
    # Returns nil on success. Raises a SudokuError on failure.
    def solve
      return if valid?

      solve_for(0)

      raise SudokuError, "failed to solve board!" unless valid?
    end

    def x_y(ptr)
      # Calculate x,y value of ptr
      x = ptr % 9
      y = ptr / 9

      [x, y]
    end

    private

    # Private: Check a position in the board Array to see whether its row
    # and column are in a valid state. 0 values are ignored.
    #
    # ptr   - Fixnum position in the array to check.
    #
    # Returns Boolean true if the row/column are valid; false if not.
    def check?(ptr)
      ptr_x, ptr_y = x_y(ptr)

      # Keep track of what we've seen
      row = {}
      col = {}
      sec = {}

      @board.each_with_index do |value, i|
        next if value == 0

        x, y = x_y(i)

        if y == ptr_y
          return false if row[value]

          row[value] = true
        end

        if x == ptr_x
          return false if col[value]

          col[value] = true
        end

        # Check the cells in the sector
        if x / 3 == ptr_x / 3 && y / 3 == ptr_y / 3
          return false if sec[value]

          sec[value] = true
        end
      end

      # These checks aren't *theoretically* necessary, given the algorithm,
      # but nice to have as a sanity check.
      if row.length == 9
        return false unless row.keys.inject(:+) == 45
      end

      if col.length == 9
        return false unless col.keys.inject(:+) == 45
      end

      true
    end

    # Private: Solve the board recursively, starting at position <ptr>.
    #
    # ptr   - Fixnum position in the array to start from.
    #
    # Returns Boolean true if the rest of the board is valid, false if the
    # current position is invalid and we've run out of numbers to try.
    def solve_for(ptr = 0)
      return true unless ptr < 81

      if @board[ptr] != 0 && check?(ptr)
        return solve_for(ptr + 1)
      end

      # Assignment instead of incrementing appears to be much faster
      NUMS.each do |n|
        @board[ptr] = n
        return false if @board[ptr] == 0

        next unless check?(ptr)
        return true if solve_for(ptr + 1)
      end
    end
  end

  class SudokuError < StandardError; end
end

board = Sudoku::CLI.main

puts "Input:" if ENV["DEBUG"]
puts board.display if ENV["DEBUG"]
board.solve
puts "Output:" if ENV["DEBUG"]
board.display

