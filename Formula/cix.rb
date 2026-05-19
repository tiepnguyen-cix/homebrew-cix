class Cix < Formula
  desc "Run and debug GitLab CI pipelines locally"
  homepage "https://github.com/tiepnguyen-cix/cix"
  version "0.1.2"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-arm64"
      sha256 "66f9a9cce058cded03a49e74ddab64c84c44a5fa9f95a4cce0c74b03240eddaa"
    else
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.1/cix-darwin-amd64"
      sha256 "a4440e7853542e985b2528f2b7a340781c40005333d652438aea1fa74efa94fe"
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
