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
        fn block, {total, idx} ->
          case block do
            {space, id} ->
              {total + Enum.sum(Enum.map(idx..(idx + space - 1), &(&1 * id))), idx + space}
            space ->
              {total, idx + space}
          end
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

  def find_suitable_space(file, block_size, block_idx) do
    file
    |> Enum.with_index()
    |> Enum.find_index(fn {space_candidate, space_idx} ->
      case space_candidate do
        {_, _} -> false
        space -> space >= block_size && block_idx > space_idx
      end
    end)
  end

  def replace_with_space(file, {block_size, _}, block_idx) do
    {new_file, _} =
      file
      |> List.replace_at(block_idx, block_size)
      |> Enum.reduce(
        {[], 0},
        fn block, {squashed_file, running_space} ->
          case block do
            {_, _} ->
              if running_space > 0 do
                {[block | [running_space | squashed_file]], 0}
              else
                {[block | squashed_file], running_space}
              end

            spaces ->
              {squashed_file, running_space + spaces}
          end
        end
      )

    Enum.reverse(new_file)
  end

  def insert_block(file, space_idx, {block_size, id}, block_idx) do
    spaces = Enum.at(file, space_idx)

    if spaces == block_size do
      List.replace_at(file, space_idx, {block_size, id})
      |> replace_with_space({block_size, id}, block_idx)
    else
      file
      |> List.replace_at(space_idx, spaces - block_size)
      |> List.insert_at(space_idx, {block_size, id})
      |> replace_with_space({block_size, id}, block_idx + 1)
    end
  end

  def find_unfragged_compact_file_checksum(blocks) do
    blocks
    |> reverse_memory_blocks()
    |> Enum.reduce(
      blocks,
      fn {block_size, id}, file ->
        block_idx = Enum.find_index(file, &(&1 == {block_size, id}))
        space_idx = find_suitable_space(file, block_size, block_idx)

        if space_idx do
          insert_block(file, space_idx, {block_size, id}, block_idx)
        else
          file
        end
      end
    )
    |> calc_checksum()
  end
end

input = File.read!("input.txt") |> Day9.parse_file()
IO.puts("Part 1: #{Day9.find_compacted_file_checksum(input)}")
IO.puts("Part 2: #{Day9.find_unfragged_compact_file_checksum(input)}")
