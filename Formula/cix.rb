class Cix < Formula
  desc "Run and debug GitLab CI pipelines locally"
  homepage "https://github.com/tiepnguyen-cix/cix"
  version "0.1.3"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-arm64"
      sha256 "eb3a9503cb6d8896c3bd40a8178e0cb137e0385d39de9b1ce45d6848b7b21ead"
    else
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-amd64"
      sha256 "58046f65a0c831b4faeaac3388f6096eb78979d7f3e79005aacebbd460fc39b5"
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
