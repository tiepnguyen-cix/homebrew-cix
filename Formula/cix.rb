class Cix < Formula
  desc "Run and debug GitLab CI pipelines locally"
  homepage "https://github.com/tiepnguyen-cix/cix"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.0/cix-darwin-arm64"
      sha256 "854f0a9eaab309ccd83a453b711e96acd77103cf7c757c0d781d20b730ccd9d5"
    else
      url "https://github.com/tiepnguyen-cix/cix/releases/download/v0.1.0/cix-darwin-amd64"
      sha256 "8a3ab253aee8936724dded6de838694d99004ea051e4c4334a80c5c5b9f7767a"
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
