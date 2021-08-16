generate-schema:
	go run schema/codegen/codegen.go --schema schema/block.elrond.avsc --out schema/schema.go
