defmodule Day9 do
  def string_list_to_int_list(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def parse_file(file) do
    file
    |> String.graphemes()
    |> string_list_to_int_list()
    |> Enum.with_index()
    |> Enum.map(fn {blocks, idx} ->
      if rem(idx, 2) == 0 do
        {blocks, round(idx / 2)}
      else
        blocks
      end
    end)
  end

  def reverse_memory_blocks(file) do
    file
    |> Enum.filter(fn block ->
      case block do
        {_, _} -> true
        _ -> false
      end
    end)
    |> Enum.reverse()
  end

  def calc_total_filled_blocks(filtered_blocks) do
    filtered_blocks
    |> Enum.reduce(0, fn {space, _}, acc -> acc + space end)
  end

  def fill_spaces(
        [blocks_head | blocks_tail],
        [{rh_space, rh_id} | reverse_tail],
        new_blocks \\ []
      ) do
    case blocks_head do
      {_, bh_id} ->
        if bh_id == rh_id do
          Enum.reverse([{rh_space, rh_id} | new_blocks])
        else
          fill_spaces(
            blocks_tail,
            [{rh_space, rh_id} | reverse_tail],
            [blocks_head | new_blocks]
          )
        end

      0 ->
        fill_spaces(
          blocks_tail,
          [{rh_space, rh_id} | reverse_tail],
          new_blocks
        )

      blank_space ->
        cond do
          rh_space > blank_space ->
            fill_spaces(
              blocks_tail,
              [{rh_space - blank_space, rh_id} | reverse_tail],
              [{blank_space, rh_id} | new_blocks]
            )

          rh_space < blank_space ->
            fill_spaces(
              [blank_space - rh_space | blocks_tail],
              reverse_tail,
              [{rh_space, rh_id} | new_blocks]
            )

          rh_space == blank_space ->
            fill_spaces(
              blocks_tail,
              reverse_tail,
              [{rh_space, rh_id} | new_blocks]
            )
        end
    end
  end

  def calc_checksum(blocks) do
    {checksum, _} =
      blocks
      |> Enum.reduce(
        {0, 0},
        fn {space, id}, {total, idx} ->
          {total + Enum.sum(Enum.map(idx..(idx + space - 1), &(&1 * id))), idx + space}
        end
      )

    checksum
  end

  def find_compacted_file_checksum(blocks) do
    reversed_blocks = reverse_memory_blocks(blocks)

    blocks
    |> fill_spaces(reversed_blocks)
    |> calc_checksum()
  end
end

input = File.read!("input.txt") |> Day9.parse_file()
IO.puts("Part 1: #{Day9.find_compacted_file_checksum(input)}")
