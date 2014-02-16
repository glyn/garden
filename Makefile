all: protocol skeleton

protocol: $(shell find protobuf/ -type f)
	mkdir -p protocol/
	rm protocol/*.pb.go
	protoc --gogo_out=protocol/ --proto_path=protobuf/ protobuf/*.proto

skeleton: warden/warden
	cd warden/warden/src && make clean all
	cp warden/warden/src/wsh/wshd root/linux/skeleton/bin
	cp warden/warden/src/wsh/wsh root/linux/skeleton/bin
	cp warden/warden/src/oom/oom root/linux/skeleton/bin
	cp warden/warden/src/iomux/iomux-spawn root/linux/skeleton/bin
	cp warden/warden/src/iomux/iomux-link root/linux/skeleton/bin
	cp warden/warden/src/repquota/repquota root/bin

warden/warden:
	git submodule update --init --recursive
