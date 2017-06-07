# Original Version and Article, Special-Thanx to :
# http://momijiame.tumblr.com/post/32458077768/spec-rpm
%define    debug_package %{nil}

Name:      goqos
Version:   0.1
Release:   %{release}
Group:     Utilities/Misc
License:   BSD
URL:       https://git.corp.yahoo.co.jp/query-engine
Summary:   TC Wrapper tool by Go
BuildArch: x86_64
Source0:   %{name}.tar.gz
Prefix:    %{_prefix}
# (only create temporary directory name, for RHEL5 compat environment)
# see : http://fedoraproject.org/wiki/Packaging:Guidelines#BuildRoot_tag
BuildRoot:  %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

Requires:      glibc
# install tc & cbq command require 3.x
Requires:      iproute
Requires:      ethtool
Requires:      /bin/sh
BuildRequires: golang

%define INSTALLDIR %{buildroot}

%description
goqos is a traffic-control tool by Go.
Original is written by perl.
https://github.com/yuokada/qos-control

%prep
%setup -q -n %{name}
# see: https://vinelinux.org/docs/vine6/making-rpm/setup-macro.html
mkdir -p $RPM_BUILD_ROOT/usr/local/{bin,man/man1}
echo $RPM_BUILD_ROOT
echo %{INSTALLDIR}

%build
make build

%install
rm   -rf %{INSTALLDIR}
mkdir -p %{buildroot}
mkdir -p %{buildroot}/etc/sysconfig/qos
mkdir -p %{buildroot}/etc/rc.d/init.d

%{__install} -D -p -m 0755 ./%{name}         %{buildroot}%{_prefix}/bin/%{name}
%{__install} -D -p -m 0755 ./bin/qos.sh      %{buildroot}/usr/local/sbin/qos.sh

%{__install} -D -p -m 0755 ./etc/rc.d/init.d/qos.init %{buildroot}/etc/init.d/qos.init
%{__install} -D -p -m 0644 ./etc/sysconfig/qos/avpkt  %{buildroot}/etc/sysconfig/qos/avpkt

# Instructions to clean out the build root.
%clean
# Avoid Disastarous Damage : http://dev.tapweb.co.jp/2010/12/273
[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf $RPM_BUILD_ROOT

%files
%defattr(0755,root,root)
%{_prefix}/bin/goqos
/usr/local/sbin/qos.sh
/etc/init.d/qos.init

%defattr(0444,root,root)
/etc/sysconfig/qos/avpkt

%post
/usr/sbin/ethtool -K eth0 tso off
/usr/sbin/ethtool -K eth0 gso off

%postun
/usr/sbin/ethtool -K eth0 tso on
/usr/sbin/ethtool -K eth0 gso on

%changelog
* Sat Apr 08 2017 yuokada
- initial release
