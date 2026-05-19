class Cix < Formula
  desc "Run and debug GitLab CI pipelines locally"
  homepage "https://github.com/tiepnguyen-cix/cix"
  version "0.1.4"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-arm64"
      sha256 "35222df7da1cc613841e373b785326d613bc2bef1bd753677d67e3b0ced51cee"
    else
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-amd64"
      sha256 "9c46a013b47a64a9f4c9d2d7f3ff2dbfe0f1b14118f1a13bd1bb0658a4d56cb3"
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
