include $(TOPDIR)/rules.mk

PKG_NAME:=dvpn-node
PKG_VERSION:=1.0.0
PKG_RELEASE:=1

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/sentinel-official/dvpn-node.git
PKG_SOURCE_VERSION:=1c5eaa789c4390e89b53968d44369c74fbbde170
PKG_MIRROR_HASH:=583e324c08a3116e33433097d0e8350505f2edfc973c327d7b051c56e3db83c0

PKG_MAINTAINER:=your_name

PKG_BUILD_DIR:=$(BUILD_DIR)/dvpn-node-$(PKG_VERSION)
PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_USE_MIPS16:=0

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

GO_PKG:=github.com/sentinel-official/dvpn-node
GO_PKG_BUILD_PKG:=github.com/sentinel-official/dvpn-node
GO_PKG_LDFLAGS:=-s -w
GO_PKG_LDFLAGS_X:= \
	github.com/cosmos/cosmos-sdk/version.Name=sentinel \
	github.com/cosmos/cosmos-sdk/version.AppName=sentinelnode \
	github.com/cosmos/cosmos-sdk/version.Version=${VERSION} \
	github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT} \
	github.com/cosmos/cosmos-sdk/version.BuildTags=${BUILD_TAGS}
GO_PKG_TAGS:= $(strip netgo,ledger,openwrt)

include $(INCLUDE_DIR)/package.mk
include $(INCLUDE_DIR)/nls.mk
include ../../lang/golang/golang-package.mk

define Package/dvpn-node
	SECTION:=net
	CATEGORY:=Network
	DEPENDS:=$(GO_ARCH_DEPENDS)
	TITLE:=dvpn-node
	URL:=https://github.com/sentinel-official/dvpn-node
endef

define Package/dvpn-node/install
	$(RM) -v $(PKG_BUILD_DIR)/services/wireguard/wireguard_linux.go

	$(call GoPackage/Package/Install/Bin,$(1))
endef

$(eval $(call GoBinPackage,dvpn-node))
$(eval $(call BuildPackage,dvpn-node))
