class Cix < Formula
  desc "Run and debug GitLab CI pipelines locally"
  homepage "https://github.com/tiepnguyen-cix/cix"
  version "0.1.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-arm64"
      sha256 "e860556e574008c24611c8042bbc9d2ca04e95830fb2965fee501c6257d18dd2"
    else
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-amd64"
      sha256 "fb0d27c7a2c427a30b36662bbcb13027ddc1dbd2d0ae2c445f4f7f6e74e86d2e"
    end
  end

  def install
    if Hardware::CPU.arm?
      bin.install "cix-darwin-arm64" => "cix"
    else
      bin.install "cix-darwin-amd64" => "cix"
    end
  end

  test do
    system "#{bin}/cix", "--version"
  end
end
