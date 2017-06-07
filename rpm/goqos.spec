# Original Version and Article, Special-Thanx to :
# http://momijiame.tumblr.com/post/32458077768/spec-rpm
%define    prefix  /
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
Prefix:    /
# (only create temporary directory name, for RHEL5 compat environment)
# see : http://fedoraproject.org/wiki/Packaging:Guidelines#BuildRoot_tag
BuildRoot: %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
Requires:      glibc
BuildRequires: golang

#%define INSTALLDIR %{buildroot}/goqos
%define INSTALLDIR %{buildroot}

%description
goqos is a traffic-control tool by Go.
Original is written by perl.
https://github.com/yuokada/qos-control

%prep
%setup -q -n %{name}
# %setup -q -n goqos
#%setup -a 0 -q
# see: https://vinelinux.org/docs/vine6/making-rpm/setup-macro.html
# mkdir -p $RPM_BUILD_ROOT/usr/local/{bin,man/man1}
# echo $RPM_BUILD_ROOT
# echo %{INSTALLDIR}
mkdir -p $RPM_BUILD_ROOT/usr/local/{bin,man/man1}
echo $RPM_BUILD_ROOT
echo %{INSTALLDIR}

%build
make build

%install
rm   -rf %{INSTALLDIR}
mkdir -p %{buildroot}
mkdir -p %{buildroot}/etc/sysconfig/qos
%{__install} -D -p -m 0755 ./%{name}  %{buildroot}%{prefix}/bin/%{name}

# Instructions to clean out the build root.
%clean
#rm -rf %{buildroot}
# Avoid Disastarous Damage : http://dev.tapweb.co.jp/2010/12/273
[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf $RPM_BUILD_ROOT

%files
%defattr(0755,root,root)
#%{prefix}/bin/goqos
%{prefix}/bin/goqos

# directory only
%dir %attr(0755,-,-) /etc/sysconfig/qos

%pre

%post

%changelog
* Sat Apr 08 2017 yuokada
- initial release
