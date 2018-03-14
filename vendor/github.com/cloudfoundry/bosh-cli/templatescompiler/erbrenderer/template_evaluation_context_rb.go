package erbrenderer

const templateEvaluationContextRb = `
# Based on common/properties/template_evaluation_context.rb
require "rubygems"
require "ostruct"
require "json"
require "erb"
require "yaml"

class Hash
  def recursive_merge!(other)
    self.merge!(other) do |_, old_value, new_value|
      if old_value.class == Hash && new_value.class == Hash
        old_value.recursive_merge!(new_value)
      else
        new_value
      end
    end
    self
  end
end

class UnknownProperty < StandardError
  attr_reader :name

  def initialize(names)
    @names = names
    super("Can't find property '#{names.join("', or '")}'")
  end
end

class ActiveElseBlock
  def initialize(template)
    @context = template
  end

  def else
    yield
  end

  def else_if_p(*names, &block)
    @context.if_p(*names, &block)
  end
end

class InactiveElseBlock
  def else; end

  def else_if_p(*names)
    InactiveElseBlock.new
  end
end

class TemplateEvaluationContext
  attr_reader :name, :index
  attr_reader :properties, :raw_properties
  attr_reader :spec

  def initialize(spec)
    File.open("/tmp/after", "w") { |f| f.write JSON.pretty_generate(spec) }

    @name = spec["job"]["name"] if spec["job"].is_a?(Hash)
    @index = spec["index"]

    if !spec['job_properties'].nil?
      properties1 = spec['job_properties']
    else
      properties1 = spec['global_properties'].recursive_merge!(spec['cluster_properties'])
    end

    properties = {}
    spec['default_properties'].each do |name, value|
      copy_property(properties, properties1, name, value)
    end

    @properties = openstruct(properties)
    @raw_properties = properties
    @spec = openstruct(spec)

    @links = spec['links'] || {}
  end

  def get_binding
    binding.taint
  end

  def p(*args)
    names = Array(args[0])

    names.each do |name|
      result = lookup_property(@raw_properties, name)
      return result unless result.nil?
    end

    return args[1] if args.length == 2
    raise UnknownProperty.new(names)
  end

  def if_p(*names)
    values = names.map do |name|
      value = lookup_property(@raw_properties, name)
      return ActiveElseBlock.new(self) if value.nil?
      value
    end

    yield *values
    InactiveElseBlock.new
  end

  def if_link(name)
    false
  end

  def link(name)
    link_spec = lookup_property(@links, name)
    raise UnknownLink.new(name) if link_spec.nil?

    if link_spec.has_key?('instances')
      return create_evaluation_link(link_spec)
    end

    raise UnknownLink.new(name)
  end

  # Run a block of code if the link given exists
  # @param [String] name of the link
  # @yield [Object] link, which is an array of instances
  def if_link(name)
    link_spec = lookup_property(@links, name)
    if link_spec.nil? || !link_spec.has_key?('instances')
      return ActiveElseBlock.new(self)
    else
      yield create_evaluation_link(link_spec)
      InactiveElseBlock.new
    end
  end

  private

  def copy_property(dst, src, name, default = nil)
    keys = name.split(".")
    src_ref = src
    dst_ref = dst

    keys.each do |key|
      src_ref = src_ref[key]
      break if src_ref.nil? # no property with this name is src
    end

    keys[0..-2].each do |key|
      dst_ref[key] ||= {}
      dst_ref = dst_ref[key]
    end

    dst_ref[keys[-1]] ||= {}
    dst_ref[keys[-1]] = src_ref.nil? ? default : src_ref
  end

  def openstruct(object)
    case object
      when Hash
        mapped = object.inject({}) { |h, (k,v)| h[k] = openstruct(v); h }
        OpenStruct.new(mapped)
      when Array
        object.map { |item| openstruct(item) }
      else
        object
    end
  end

  def lookup_property(collection, name)
    keys = name.split(".")
    ref = collection

    keys.each do |key|
      ref = ref[key]
      return nil if ref.nil?
    end

    ref
  end

  def create_evaluation_link(link_spec)
    link_instances = link_spec['instances'].map do |instance_link_spec|
      EvaluationLinkInstance.new(
        instance_link_spec['name'],
        instance_link_spec['index'],
        instance_link_spec['id'],
        instance_link_spec['az'],
        instance_link_spec['address'],
        instance_link_spec['properties'],
        instance_link_spec['bootstrap'],
      )
    end

    # if link_spec.has_key?('address')
    #   encoder_to_inject = ManualLinkDnsEncoder.new(link_spec['address'])
    # else
    #   encoder_to_inject = @dns_encoder
    # end

    return EvaluationLink.new(
      link_instances,
      link_spec['properties'],
      link_spec['instance_group'],
      link_spec['default_network'],
      link_spec['deployment_name'],
      link_spec['domain'],
      # encoder_to_inject,
    )
  end
end

class EvaluationLinkInstance
  attr_reader :name
  attr_reader :index
  attr_reader :id
  attr_reader :az
  attr_reader :address
  attr_reader :properties
  attr_reader :bootstrap

  def initialize(name, index, id, az, address, properties, bootstrap)
    @name = name
    @index = index
    @id = id
    @az = az
    @address = address
    @properties = properties
    @bootstrap = bootstrap
  end

  def p(*args)
    names = Array(args[0])

    names.each do |name|
      result = lookup_property(@properties, name)
      return result unless result.nil?
    end

    return args[1] if args.length == 2
    raise UnknownProperty.new(names)
  end

  def if_p(*names)
    values = names.map do |name|
      value = lookup_property(@properties, name)
      return ActiveElseBlock.new(self) if value.nil?
      value
    end

    yield *values
    InactiveElseBlock.new
  end

  private

  def lookup_property(collection, name)
    keys = name.split(".")
    ref = collection

    keys.each do |key|
      ref = ref[key]
      return nil if ref.nil?
    end

    ref
  end
end

class EvaluationLink
  attr_reader :instances
  attr_reader :properties

  def initialize(instances, properties, instance_group, default_network, deployment_name, root_domain) # , dns_encoder)
    @instances = instances
    @properties = properties
    @instance_group = instance_group
    @default_network = default_network
    @deployment_name = deployment_name
    @root_domain = root_domain
    # @dns_encoder = dns_encoder
  end

  def p(*args)
    names = Array(args[0])

    names.each do |name|
      result = lookup_property(@properties, name)
      return result unless result.nil?
    end

    return args[1] if args.length == 2
    raise UnknownProperty.new(names)
  end

  def if_p(*names)
    values = names.map do |name|
      value = lookup_property(@properties, name)
      return ActiveElseBlock.new(self) if value.nil?
      value
    end

    yield *values
    InactiveElseBlock.new
  end

  def address(criteria = {})
    raise NotImplementedError.new('link.address requires bosh director') # if @dns_encoder.nil?

    full_criteria = criteria.merge(
      instance_group: @instance_group,
      default_network: @default_network,
      deployment_name: @deployment_name,
      root_domain: @root_domain,
    )

    @dns_encoder.encode_query(full_criteria)
  end

  private

  def lookup_property(collection, name)
    keys = name.split(".")
    ref = collection

    keys.each do |key|
      ref = ref[key]
      return nil if ref.nil?
    end

    ref
  end
end

# todo do not use JSON in releases
class << JSON
  alias dump_array_or_hash dump

  def dump(*args)
    arg = args[0]
    if arg.is_a?(String) || arg.is_a?(Numeric)
      arg.inspect
    else
      dump_array_or_hash(*args)
    end
  end
end

class ERBRenderer
  def initialize(context)
    @context = context
  end

  def render(src_path, dst_path)
    erb = ERB.new(File.read(src_path), safe_level = nil, trim_mode = "-")
    erb.filename = src_path

    File.open(dst_path, "w") do |f|
      f.write(erb.result(@context.get_binding))
    end

  rescue Exception => e
    name = "#{@context.name}/#{@context.index}"

    line_i = e.backtrace.index { |l| l.include?(erb.filename) }
    line_num = line_i ? e.backtrace[line_i].split(':')[1] : "unknown"
    location = "(line #{line_num}: #{e.inspect})"

    raise("Error filling in template '#{src_path}' for #{name} #{location}")
  end
end

if $0 == __FILE__
  context_path, src_path, dst_path = *ARGV

  context_hash = JSON.load(File.read(context_path))
  context = TemplateEvaluationContext.new(context_hash)

  renderer = ERBRenderer.new(context)
  renderer.render(src_path, dst_path)
end
`
