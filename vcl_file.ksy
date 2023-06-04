meta:
  id: vcl_file
  endian: le
  seq:
    - id: signature
      type: u4
      enum: signature_enum
    - id: version
      type: u4
    - id: num_textures
      type: u2
    - id: texture_headers
      type: texture_header
      repeat: expr
      repeat-expr: num_textures
    - id: palette
      type: u1
      repeat: 256
    - id: textures
      type: texture_data
      repeat: expr
      repeat-expr: num_textures
    - id: num_sounds
      type: u2
    - id: sound_headers
      type: sound_header
      repeat: expr
      repeat-expr: num_sounds
    - id: sounds
      type: sound_data
      repeat: expr
      repeat-expr: num_sounds

types:
  signature_enum:
    0x4C434556: LE
    0x5643454C: VCEL

  texture_header:
    seq:
      - id: width
        type: u2
      - id: height
        type: u2
      - id: offset
        type: u4

  texture_data:
    seq:
      - id: width
        type: u2
      - id: height
        type: u2
      - id: data
        size: calculated_texture_size

  sound_header:
    seq:
      - id: sample_rate
        type: u2
      - id: num_channels
        type: u1
      - id: bits_per_sample
        type: u1
      - id: data_offset
        type: u4

  sound_data:
    seq:
      - id: size
        type: u4
      - id: data
        size: size

  calculated_texture_size: ((width + 1) / 2) * ((height + 1) / 2)

