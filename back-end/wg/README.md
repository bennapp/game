# World Generation

Contains services to detect if portion of the world has been generated world_generation_detection.go

run_world_generation.go listens to channel to generate coord regions and then spins up a go routine to generate a region of world

generate_world.go generates the terrain for a region of the world
