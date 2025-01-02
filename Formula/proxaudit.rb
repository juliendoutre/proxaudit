# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Proxaudit < Formula
  desc ""
  homepage ""
  version "1.9.0"

  depends_on "mkcert" => "1.4.4"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/juliendoutre/proxaudit/releases/download/v1.9.0/proxaudit_Darwin_x86_64.tar.gz"
      sha256 "87e87adea2c464d3fa9ad0c3141d98ef5b83a54401ae391353157932e1eb5fb4"

      def install
        bin.install "proxaudit"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/juliendoutre/proxaudit/releases/download/v1.9.0/proxaudit_Darwin_arm64.tar.gz"
      sha256 "a0c8f0e133d21745dc786ce9ef7382f04575a33b0886d923cd3864cbfddc50b9"

      def install
        bin.install "proxaudit"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/juliendoutre/proxaudit/releases/download/v1.9.0/proxaudit_Linux_x86_64.tar.gz"
        sha256 "06c041a430584ba39b96647fe202ab0e3897ee2105e3e1f9c96d20d91a548214"

        def install
          bin.install "proxaudit"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/juliendoutre/proxaudit/releases/download/v1.9.0/proxaudit_Linux_arm64.tar.gz"
        sha256 "dde909e501611ce577a63d02d1f7664cc0d5214c1af98e0296d1723384fcd32b"

        def install
          bin.install "proxaudit"
        end
      end
    end
  end
end
